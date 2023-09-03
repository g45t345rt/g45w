Function Initialize() Uint64
1 IF EXISTS("owner") THEN GOTO 12
2 STORE("owner", SIGNER())
3 STORE("totalDonated", 0)
4 STORE("totalDonations", 0)
5 STORE("highestDonation", 0)
6 STORE("highestDonationAddr", "")
7 STORE("highestDonationTimestamp", 0)
8 STORE("lastDonation", 0)
9 STORE("lastDonationAddr", "")
10 STORE("lastDonationTimestamp", 0)
11 STORE("d_", 0)
12 RETURN 0
End Function

Function Donate() Uint64
1 DIM signer as String
2 DIM donation as Uint64
3 LET signer = ADDRESS_STRING(SIGNER())
4 LET donation = DEROVALUE()
5 IF donation == 0 THEN GOTO 19
6 STORE("totalDonated", LOAD("totalDonated") + donation)
7 STORE("totalDonations", LOAD("totalDonations") + 1)
8 IF EXISTS("d_" + signer) THEN GOTO 10
9 STORE("d_" + signer, 0)
10 STORE("d_" + signer, LOAD("d_" + signer) + donation)
11 IF donation < LOAD("highestDonation") THEN GOTO 15
12 STORE("highestDonation", donation)
13 STORE("highestDonationAddr", signer)
14 STORE("highestDonationTimestamp", BLOCK_TIMESTAMP())
15 STORE("lastDonation", donation)
16 STORE("lastDonationAddr", signer)
17 STORE("lastDonationTimestamp", BLOCK_TIMESTAMP())
18 SEND_DERO_TO_ADDRESS(LOAD("owner"), donation)
19 RETURN 0
End Function

Function TransferOwnership(newOwner String) Uint64
1 IF LOAD("owner") != SIGNER() THEN GOTO 3
2 STORE("tempOwner", ADDRESS_RAW(newOwner))
3 RETURN 0
End Function

Function CancelTransferOwnership() Uint64
1 IF LOAD("owner") != SIGNER() THEN GOTO 3
2 DELETE("tempOwner")
3 RETURN 0
End Function

Function ClaimOwnership() Uint64
1 IF LOAD("tempOwner") != SIGNER() THEN GOTO 4
2 STORE("owner", SIGNER())
3 DELETE("tempOwner")
4 RETURN 0
End Function
