/* Credit to Pieswap */

package fastreg

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/deroproject/derohe/cryptography/crypto"
	"github.com/deroproject/derohe/transaction"
	"golang.org/x/crypto/sha3"
)

const N = 32 * 1024

type result_t struct {
	txn    *transaction.Transaction
	secret *big.Int
}

type pt_t struct {
	x, y, secret *big.Int
}

func mySetBig(s string) *big.Int {
	i, _ := new(big.Int).SetString(s, 0)
	return i
}

// Dero's G1
var GXMONT = mySetBig("0x26b1948aba4465c168faeb79586c8022afdf98edf532bd9d62ce90cd019391e5")
var GYMONT = mySetBig("0x24c3da4404b3a717f30e6db01e4ef1d4cea805087c00ad93893ab6f1e1aa1f4e")
var G = &pt_t{GXMONT, GYMONT, ZERO}

// Dero's field
var P = mySetBig("21888242871839275222246405745257275088696311157297823662689037894645226208583")
var NP = mySetBig("0xf57a22b791888c6bd8afcbd01833da809ede7d651eca6ac987d20782e4866389")
var R2 = mySetBig("0x06d89f71cab8351f47ab1eff0a417ff6b5e71911d44501fbf32cfc5b538afa89")
var O = mySetBig("0x30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f0000001")
var NO = mySetBig("0x73f82f1d0d8341b2e39a9828990623916586864b4c6911b3c2e1f593efffffff")
var OR = mySetBig("0x54a47462623a04a7ab074a58680730147144852009e880ae620703a6be1de925")
var OR2 = mySetBig("0x0216d0b17f4e44a58c49833d53bb808553fe3ab1e35c59e31bb8e645ae216da7")

var ONE = mySetBig("1")
var ZERO = mySetBig("0")
var inc = &pt_t{GXMONT, GYMONT, ONE}

func randPt(bits int64) *pt_t {
	rnd := big.NewInt(0).Exp(big.NewInt(2), big.NewInt(bits), nil)
	rnd, _ = rand.Int(rand.Reader, rnd)

	p := newPt()
	curveMul(p, G, rnd)
	p.secret.Set(rnd)

	return p
}

func newPt() *pt_t {
	p := pt_t{new(big.Int).Set(ZERO), new(big.Int).Set(ZERO), new(big.Int).Set(ZERO)}

	return &p
}

func (a *pt_t) set(b *pt_t) {
	a.x = b.x
	a.y = b.y
}

func (a *pt_t) setInfinity() {
	a.x.Set(ZERO)
	a.y.Set(ZERO)
}

func (p *pt_t) String() string {
	return fmt.Sprintf("bn256.G1(%064x, %064x)", p.x, p.y)
}

// https://www.nayuki.io/page/barrett-reduction-algorithm
func barretReductionO(a *big.Int) {
	t := new(big.Int).Mul(a, OR)
	t.Rsh(t, 508)
	t.Mul(t, O)
	a.Sub(a, t)

	if a.Cmp(O) > 0 {
		a.Sub(a, O)
	}
}

func modMulP(p, a, b *big.Int) {
	modMul(p, a, b, P, NP)
}

func modMulO(p, a, b *big.Int) {
	modMul(p, a, b, O, NO)
}

func modMul(p, a, b, M, NM *big.Int) {
	var T = new(big.Int)
	var m = new(big.Int)
	var t = new(big.Int)
	var abs []big.Word

	if a.Cmp(ZERO) == 0 || b.Cmp(ZERO) == 0 {
		p.Set(ZERO)
		return
	}

	T.Mul(a, b)
	abs = T.Bits()

	m.Mul(new(big.Int).SetBits(abs[0:4]), NM)
	abs = m.Bits()

	t.Mul(new(big.Int).SetBits(abs[0:4]), M)

	p.Add(T, t)

	abs = p.Bits()
	p.SetBits(abs[4:8])
	if p.Cmp(M) > 0 {
		p.Sub(p, M)
	} // don't know why this is needed sometimes
}

// https://crypto.stackexchange.com/a/54623/555
func myInvert(a, p *big.Int) (r *big.Int) {
	u := new(big.Int).Set(a)
	v := new(big.Int).Set(p)
	x2 := new(big.Int).Set(ZERO)
	x1 := new(big.Int).Set(p)
	x1.Sub(x1, ONE)

	if u.Bit(0) == 0 {
		u.Add(a, p)
	}

	for v.Cmp(ONE) != 0 {
		for v.Cmp(u) < 0 {
			u.Sub(u, v)
			x1.Add(x1, x2)

			for u.Bit(0) == 0 {
				if x1.Bit(0) == 1 {
					x1.Add(x1, p)
				}
				u.Rsh(u, 1)
				x1.Rsh(x1, 1)
			}
		}

		v.Sub(v, u)
		x2.Add(x2, x1)

		for v.Bit(0) == 0 {
			if x2.Bit(0) == 1 {
				x2.Add(x2, p)
			}
			v.Rsh(v, 1)
			x2.Rsh(x2, 1)
		}
	}

	r = x2

	return
}

func montEncode(r, a *big.Int) {
	modMulP(r, a, R2)
}

func montDecode(r, a *big.Int) {
	modMulP(r, a, ONE)
}

func curveAddG(c, a *pt_t) {
	curveAdd(c, a, G)
}

// https://en.wikipedia.org/wiki/Elliptic_curve_point_multiplication
func curveDouble(c, a *pt_t) {
	var t1 = new(big.Int)
	var t2 = new(big.Int)
	var inv = new(big.Int)
	var s = new(big.Int)
	var s2 = new(big.Int)
	var rx = new(big.Int)
	var ry = new(big.Int)

	modMulP(t1, a.y, big.NewInt(2))
	modMulP(t2, a.x, a.x)
	modMulP(t2, t2, big.NewInt(3))

	montDecode(inv, t1)
	inv.ModInverse(inv, P)
	montEncode(inv, inv)

	modMulP(s, t2, inv)
	modMulP(s2, s, s)

	rx.Sub(s2, a.x)
	if rx.Sign() < 0 {
		rx.Add(rx, P)
	}

	rx.Sub(rx, a.x)
	if rx.Sign() < 0 {
		rx.Add(rx, P)
	}

	ry.Sub(a.x, rx)
	if ry.Sign() < 0 {
		ry.Add(ry, P)
	}

	modMulP(ry, ry, s)

	ry.Sub(ry, a.y)
	if ry.Sign() < 0 {
		ry.Add(ry, P)
	}

	c.x = rx
	c.y = ry
}

func curveAdd(c, a, b *pt_t) {
	if a.x.Cmp(ZERO) == 0 && a.y.Cmp(ZERO) == 0 {
		c.x = b.x
		c.y = b.y
		return
	}
	if b.x.Cmp(ZERO) == 0 && b.y.Cmp(ZERO) == 0 {
		c.x = a.x
		c.y = a.y
		return
	}

	if a.x == b.x && a.y == b.y {
		curveDouble(c, a)
		return
	}

	var t1 = new(big.Int)
	var t2 = new(big.Int)
	var inv = new(big.Int)
	var s = new(big.Int)
	var s2 = new(big.Int)
	var rx = new(big.Int)
	var ry = new(big.Int)

	t1.Sub(b.x, a.x)
	if t1.Sign() < 0 {
		t1.Add(t1, P)
	}

	t2.Sub(b.y, a.y)
	if t2.Sign() < 0 {
		t2.Add(t2, P)
	}

	montDecode(inv, t1)
	//	inv = myInvert(inv, P)
	inv.ModInverse(inv, P)
	montEncode(inv, inv)

	modMulP(s, t2, inv)
	modMulP(s2, s, s)

	rx.Sub(s2, a.x)
	if rx.Sign() < 0 {
		rx.Add(rx, P)
	}

	rx.Sub(rx, b.x)
	if rx.Sign() < 0 {
		rx.Add(rx, P)
	}

	ry.Sub(a.x, rx)
	if ry.Sign() < 0 {
		ry.Add(ry, P)
	}

	modMulP(ry, ry, s)

	ry.Sub(ry, a.y)
	for ry.Sign() < 0 {
		ry.Add(ry, P)
	}

	c.x = rx
	c.y = ry
}

func curveMul(r, a *pt_t, s *big.Int) {
	res := newPt()
	tmp := newPt()
	tmp.set(a)

	for i := 0; i < s.BitLen(); i++ {
		if s.Bit(i) == 1 {
			curveAdd(res, res, tmp)
		}
		curveDouble(tmp, tmp)
	}

	r.set(res)
}

func pointFactory(p *pt_t) *pt_t {
	curveAdd(p, p, inc)
	p.secret.Add(p.secret, inc.secret)

	return p
}

func listInit(ctx *pt_t) []*pt_t {
	pList := make([]*pt_t, N)

	// initial point
	p := pointFactory(ctx)

	// make point list
	for i := 0; i < N; i++ {
		pList[i] = newPt()
		montDecode(pList[i].x, p.x)
		montDecode(pList[i].y, p.y)
		pList[i].secret.Set(p.secret)

		p = pointFactory(ctx)
	}

	return pList
}

func listNext(pList []*pt_t) {
	var scratch = make([]big.Int, N)
	var accum = new(big.Int).Set(ONE)
	var accum_inv = new(big.Int)
	var inv = new(big.Int)
	var t1 = new(big.Int)
	var t2 = new(big.Int)
	var rx = new(big.Int)
	var ry = new(big.Int)
	var s = new(big.Int)
	var s2 = new(big.Int)

	for i := 0; i < N; i++ {
		t1.Sub(G.x, pList[i].x)
		if t1.Sign() < 0 {
			t1.Add(t1, P)
		}

		scratch[i].Set(accum)
		modMulP(accum, accum, t1)
	}

	montDecode(accum, accum)
	accum.ModInverse(accum, P)
	//	accum = myInvert(accum, P)
	montEncode(accum_inv, accum)

	for i := N - 1; i >= 0; i-- {
		t1.Sub(G.x, pList[i].x)
		if t1.Sign() < 0 {
			t1.Add(t1, P)
		}

		modMulP(inv, accum_inv, &scratch[i])
		modMulP(accum_inv, accum_inv, t1)

		t2.Sub(G.y, pList[i].y)
		if t2.Sign() < 0 {
			t2.Add(t2, P)
		}

		modMulP(s, t2, inv)
		modMulP(s2, s, s)

		rx.Sub(s2, pList[i].x)
		if rx.Sign() < 0 {
			rx.Add(rx, P)
		}

		rx.Sub(rx, G.x)
		if rx.Sign() < 0 {
			rx.Add(rx, P)
		}

		ry.Sub(pList[i].x, rx)
		if ry.Sign() < 0 {
			ry.Add(ry, P)
		}

		modMulP(ry, ry, s)

		ry.Sub(ry, pList[i].y)
		for ry.Sign() < 0 {
			ry.Add(ry, P)
		}

		/* validation
		   var c = newPt()
		   ns := new(big.Int).Add(pList[i].secret, ONE)
		   curveMul(c, G, ns)
		   //curveAdd(c, G, pList[i])
		   if rx.Cmp(c.x) != 0 {
		   	fmt.Printf("X = %064x\n", pList[i].x)
		   	fmt.Printf("Y = %064x\n", pList[i].y)
		   	fmt.Printf("cx = %064x\n", c.x)
		   	fmt.Printf("rx = %064x\n", rx)
		   	fmt.Printf("cy = %064x\n", c.y)
		   	fmt.Printf("ry = %064x\n", ry)
		   	fmt.Printf("ERROR [%d]\n", i)
		   }
		*/
		pList[i].x.Set(rx)
		pList[i].y.Set(ry)
		pList[i].secret.Add(pList[i].secret, ONE)
	}
}

func Search(onHash func(), onFound func(tx *transaction.Transaction, secret *big.Int)) {
	ctx := randPt(255)
	pList := listInit(ctx)

	tmpPoint := randPt(255)

	j := 0
	for {
		tp := newPt()
		montDecode(tp.x, tmpPoint.x)
		montDecode(tp.y, tmpPoint.y)
		tp.secret.Set(tmpPoint.secret)

		for i := 0; i < N; i++ {
			txn := getRegistrationTX(pList[i], tp)
			onHash()

			hash := GetHash(txn)
			if hash[0] == 0 && hash[1] == 0 && hash[2] == 0 {
				if txn.IsRegistrationValid() {
					fmt.Println("Found valid registration tx")
					onFound(txn, pList[i].secret)
					pList[i] = pointFactory(ctx)
					break
				}
				fmt.Println("Found registration tx but invalid. Let's continue.")
			}
			j++
		}

		tmpPoint = pointFactory(tmpPoint)
		time.Sleep(100 * time.Millisecond)
	}
}

func getRegistrationTX(p, tp *pt_t) *transaction.Transaction {
	var tx transaction.Transaction
	tx.Version = 1
	tx.TransactionType = transaction.REGISTRATION

	addr := compress(p)
	copy(tx.MinerAddress[:], addr[:])

	c, s := sign(p, tp)
	c.FillBytes(tx.C[:])
	s.FillBytes(tx.S[:])

	/*
		if !tx.IsRegistrationValid() {
			panic("registration tx could not be generated. something failed.")
		} else {
			fmt.Println("Reg OK")
		}
	*/

	return &tx
}

func sign(p, tp *pt_t) (c, s *big.Int) {
	serialize := []byte(p.String() + tp.String()) // 280 bytes
	c = HashtoNumber(serialize)
	//	barretReductionO(c)
	for true {
		t := new(big.Int).Sub(c, O)
		if t.Sign() > 0 {
			c = t
		} else {
			break
		}
	}

	s = new(big.Int).Mul(c, p.secret)
	barretReductionO(s)

	//	s = new(big.Int)
	//	modMulO(s, c, OR2)
	//	modMulO(s, s, p.secret)

	s.Add(s, tp.secret)

	t := new(big.Int).Sub(s, O)
	if t.Sign() > 0 {
		s = t
	}
	//	barretReductionO(s)

	return
}

func compress(p *pt_t) []byte {
	b := make([]byte, 32)
	p.x.FillBytes(b)
	b = append(b, 0x00)

	y2 := new(big.Int).Sub(P, p.y)
	if p.y.Cmp(y2) >= 0 {
		b[32] = 0x01
	}

	return b
}

func HashtoNumber(input []byte) *big.Int {
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(input)

	hash := hasher.Sum(nil)
	return new(big.Int).SetBytes(hash[:])
}

const HashLength = 32

type Hash [HashLength]byte

func GetHash(tx *transaction.Transaction) (result Hash) {
	switch tx.Version {
	case 1:
		result = Hash(crypto.Keccak256(tx.SerializeCoreStatement())) // 101 bytes
	default:
		panic("Transaction version unknown")

	}
	return
}
