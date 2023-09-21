import { glob, keyframes, styled } from "goober"

glob`
  @font-face {
    font-family: Roboto;
    src: url(./fonts/Roboto-Regular.ttf) format('truetype');
    font-display: fallback;
  }

  html, body {
    font-family: Roboto;

    /* thanks to https://www.magicpattern.design/tools/css-backgrounds */
    background-color: #1a1a1a;
    background-image: radial-gradient(#2c2c2c 1.35px, #1a1a1a 1.35px);
    background-size: 27px 27px;
    margin: 0 1em;
  }
`

const rotate = keyframes`
    from, to {
        transform: rotate(0deg);
    }

    100% {
        transform: rotate(360deg);
    }
`;

export const RotateIcon = styled("span")`
  animation: ${rotate} 1s ease-in-out infinite;
`

export const Box = styled("div")`
  max-width: 350px;
  background-color: white;
  border-radius: 10px;
  padding: 2em;
  margin: 5em auto;
  position: relative;
`

export const BoxTitle = styled("div")`
  font-size: 34px;
  font-weight: bold;
  text-align: center;
  margin: 10px 0;
`

export const BoxDescription = styled("div")`
  font-size: 20px;
  text-align: center;
  margin: 10px 0 20px 0;
`

export const BoxSubmitButton = styled("button")`
  border-radius: 10px;
  padding: 0.7em;
  font-size: 20px;
  font-weight: bold;
  border: none;
  width: 100%;
  cursor: pointer;
  margin-top: 20px;
  background-color: #f0f0f0;
  display: flex;
  align-items: center;
  gap: .5em;
  justify-content: center;
`

export const ErrDiv = styled("div")`
  border: 3px solid red;
  border-radius: 10px;
  padding: 1em;
  color: gray;
  word-break: break-all;
`

export const WarningDiv = styled("div")`
  border: 3px solid yellow;
  border-radius: 10px;
  padding: 1em;
  color: gray;
`

export const SuccessDiv = styled("div")`
  border: 3px solid green;
  border-radius: 10px;
  padding: 1em;
  color: gray;
`

export const BoxItemTitle = styled("div")`
  margin-bottom: 5px;
`

export const BoxItemValue = styled("div")`
  margin-bottom: 10px;
  color: gray;
  word-break: break-all;
  font-size: 14px;
`

export const BridgeTransfer = styled("div")`
  padding: 1em;
  border-radius: 10px;
  margin: 20px 0;
  font-size: 26px;
  border: 3px solid #f0f0f0;
  text-align: center;
  font-weight: bold;
`

export const MetaMaskLogoWrap = styled("div")`
  padding: 15px;
  left: 50%;
  margin-left: -40px;
  top: -40px;
  width: 50px;
  height: 50px;
  border-radius: 50%;
  background-color: white;
  position: absolute;
`

export const MetaMaskLogo = styled("div")`
  background-image: url(./images/metamask.png);
  background-size: contain;
  background-repeat: no-repeat;
  width: 100%;
  height: 100%;
`