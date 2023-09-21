import { h, render } from "preact"
import { useCallback, useEffect, useState } from "preact/hooks"
import * as ethers from "ethers"
import { formatEther, formatUnits, parseUnits } from "ethers"
import { setup } from 'goober'
import detectEthProvider from '@metamask/detect-provider'
import * as base64 from 'js-base64'

import 'reset-css'
import 'material-icons/iconfont/round.css'

import { bridgeAbi, tokenAbi, bridgeAddress } from "./bridge"
import {
  Box, BoxDescription, BoxItemTitle, BoxItemValue,
  BoxSubmitButton, BoxTitle, BridgeTransfer, ErrDiv, MetaMaskLogo,
  MetaMaskLogoWrap, SuccessDiv, WarningDiv, RotateIcon
} from "./style"

setup(h)


const BoxItem = (props) => {
  const { title, value } = props

  var render = value
  if (typeof render === "function") {
    render = value()
  }

  return <div>
    <BoxItemTitle>{title}</BoxItemTitle>
    <BoxItemValue>{render}</BoxItemValue>
  </div>
}

const App = () => {
  const [loaded, setLoaded] = useState()
  const [err, setErr] = useState()
  const [chainId, setChainId] = useState()
  const [account, setAccount] = useState()
  const [bridgeData, setBridgeData] = useState()
  const [provider, setProvider] = useState()
  const [bridging, setBridging] = useState()
  const [bridgeFee, setBridgeFee] = useState()
  const [bridgeHash, setBridgeHash] = useState()

  // Initialize provider and querystring data
  useEffect(async () => {
    const isBridgeDataValid = (bridgeData) => {
      if (!bridgeData) return false
      if (!bridgeData.walletAddress || !bridgeData.symbol || !bridgeData.amount) return false

      return true
    }

    try {
      const queryString = (new URL(location)).searchParams
      const base64Data = base64.decode(queryString.get("data"))
      const queryData = JSON.parse(base64Data)

      if (!isBridgeDataValid(queryData)) {
        throw ""
      }

      setBridgeData(queryData)
    } catch (err) {
      setErr(new Error("Missing or invalid bridge data. Close and request a new instance from the application."))
      return
    }

    const provider = await detectEthProvider()
    if (!provider) {
      setErr(new Error("Metamask is not available or multiple wallets installed?"))
      return
    }

    setProvider(new ethers.BrowserProvider(window.ethereum))
    setLoaded(true)
  }, [])

  // Handle chainId
  useEffect(async () => {
    if (!loaded) return

    const chainId = await window.ethereum.request({ method: "eth_chainId" })
    setChainId(chainId)

    const handleChainChanged = (chainId) => {
      setChainId(chainId)
    }

    window.ethereum.on("chainChanged", handleChainChanged)
  }, [loaded])

  // Handle accounts
  useEffect(async () => {
    if (!loaded) return

    try {
      const accounts = await window.ethereum.request({ method: "eth_requestAccounts" })
      setAccount(accounts[0])
    } catch (err) {
      setErr(err)
    }

    await window.ethereum.request({ method: "eth_accounts" })

    const handleAccountsChanged = (accounts) => {
      setAccount(accounts[0])
    }

    window.ethereum.on("accountsChanged", handleAccountsChanged)
  }, [loaded])

  // Handle bridge fee
  useEffect(async () => {
    if (!provider && !bridgeData) return

    const bridgeContract = new ethers.Contract(bridgeAddress, bridgeAbi, provider)
    const bridgeFee = await bridgeContract.bridgeFee()
    setBridgeFee(bridgeFee)
  }, [provider])

  // Call bridge smart contract
  const bridgeIn = useCallback(async () => {
    if (!provider && !bridgeData) return

    try {
      setErr(null)
      setBridging(true)

      const { walletAddress, amount, symbol } = bridgeData

      const signer = await provider.getSigner()
      const userBridgeContract = new ethers.Contract(bridgeAddress, bridgeAbi, signer)

      const bridgeContract = new ethers.Contract(bridgeAddress, bridgeAbi, provider)
      const tokenAddress = await bridgeContract.registeredSymbol(symbol)

      const tokenContract = new ethers.Contract(tokenAddress, tokenAbi, provider);
      const tokenDecimals = await tokenContract.decimals()

      const tokenAllowance = await tokenContract.allowance(account, bridgeAddress)
      const amountToBridge = parseUnits(amount.toString(), tokenDecimals)

      const userTokenContract = new ethers.Contract(tokenAddress, tokenAbi, signer)

      if (amountToBridge > tokenAllowance) {
        await userTokenContract.approve(bridgeAddress, amountToBridge)
      }

      const options = { value: bridgeFee }
      const tx = await userBridgeContract.bridgeETH2DERO(tokenAddress, walletAddress, amountToBridge, options)
      setBridgeHash(tx.hash)
    } catch (err) {
      if (err.info && err.info.error) {
        setErr(err.info.error)
      } else {
        setErr(err)
      }
    }
    setBridging(false)
  }, [provider, bridgeFee, bridgeData])

  return <div>
    <Box>
      <MetaMaskLogoWrap>
        <MetaMaskLogo />
      </MetaMaskLogoWrap>
      <BoxTitle>G45W</BoxTitle>
      <BoxDescription>Bridge Ethereum to Dero Stargate with Metamask.</BoxDescription>
      {(() => {
        if (err) {
          return <div>
            <ErrDiv>
              {err.message}
            </ErrDiv>
            <BoxSubmitButton onClick={() => location.reload()}>
              RELOAD
            </BoxSubmitButton>
          </div>
        }

        if (!loaded) return

        if (!account) {
          return <WarningDiv>
            Waiting for user to connect Metamask...
          </WarningDiv>
        }

        if (bridgeHash) {
          return <SuccessDiv>
            Bridging successful. Your wrapped tokens will appear in a couple of minutes.<br /><br />
            <a href={`https://etherscan.io/tx/${bridgeHash}`} target="_blank" style="word-break:break-all;">{bridgeHash}</a>
          </SuccessDiv>
        }

        return <div>
          <BoxItem title="Bridge Contract" value={bridgeAddress} />
          <BoxItem title="Ethereum Address" value={account} />
          <BoxItem title="Dero Address" value={bridgeData.walletAddress} />
          <BridgeTransfer>
            {bridgeData.amount + " " + bridgeData.symbol}
          </BridgeTransfer>
          <BoxItem title="Bridge In Fee" value={() => {
            return formatEther((bridgeFee || 0).toString()) + " ETH"
          }} />
          <div>
            <BoxSubmitButton onClick={bridgeIn} disabled={bridging}>
              {bridging && <RotateIcon className="material-icons-round">refresh</RotateIcon>}
              {bridging ? "BRIDGING" : "BRIDGE IN"}
            </BoxSubmitButton>
          </div>
        </div>
      })()}
    </Box>
  </div>
}

render(<App />, document.getElementById("app"))