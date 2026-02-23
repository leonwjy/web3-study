const { ethers, upgrades } = require("hardhat")
const { Side, SaleKind } = require("../test/common")
const { toBn } = require("evm-bn")

/**  * 2024/12/22 in sepolia testnet
 * 
 * 
 * deployer:  0xBb03D7e8380a60665710582E42668cFceF85f457
    esVault contract deployed to: 0x9f09e574f507BA195b2138474DDA451B4032b96D
    0xc17268561c1E5b21C9224809DF6456Fa86FD87A1  esVault getImplementationAddress
    0xd7733630336a51D14ACadbB8289a1124B4B92b5f  esVault getAdminAddress
  esDex contract deployed to: 0x0Bbe8D3974A007A0Ae61F19BeCB059F00d95532b
    0xA35eC7171a97C262d01731A15EC4a7d095ECFcc6  esDex getImplementationAddress
    0xd7733630336a51D14ACadbB8289a1124B4B92b5f  esDex getAdminAddress
  esVault setOrderBook tx: 0x278370f06a8d28387e9dab0cf018a678d03c888575c0a5823465a8e3ba22f4a3
 */

const esDex_name = "EasySwapOrderBook";
const esDex_address = "0x0Bbe8D3974A007A0Ae61F19BeCB059F00d95532b"

const esVault_name = "EasySwapVault";
const esVault_address = "0x9f09e574f507BA195b2138474DDA451B4032b96D"

const erc721_name = "TestERC721"
const erc721_address = "0x5AD7186d8cf84091b4823DB49198b03207E5247E"

let esDex, esVault, testERC721
let deployer
async function main() {
    [deployer, trader] = await ethers.getSigners()
    console.log("deployer: ", deployer.address)
    console.log("trader: ", trader.address)

    esDex = await (
        await ethers.getContractFactory(esDex_name)
    ).attach(esDex_address)

    esVault = await (
        await ethers.getContractFactory(esVault_name)
    ).attach(esVault_address)

    testERC721 = await (
        await ethers.getContractFactory(erc721_name)
    ).attach(erc721_address)


    // 1. setApprovalForAll
    await approvalForVault();

    // 2. make order
    // await testMakeOrder();

    // for (let i = 0; i < 20; i++) {
    //     await testMakeOrder(i);
    // }

    // 3. cancel order
    // let orderKeys = [];
    // await testCancelOrder(orderKeys);

    // let orderKeys1 = ["0x0c30bfe7f32dcac09941839b595a8c075d70588d4eace41f54b97cee02577df4"]
    // let orderKeys2 = ["0x5472a27cef6fba8a86d80adafe94e07c5867b3339dddac78ca9ffb41a3fdfe6a",
    //     "0x0c30bfe7f32dcac09941839b595a8c075d70588d4eace41f54b97cee02577df4"]

    // await testCancelOrder(orderKeys1);
    // await testCancelOrder(orderKeys2);


    // 4. match order 
    // await testMatchOrder();

    let orderKeys = ["0xfc20766a6184743924adb56d286d1f8de765a438642dddc6f1011e22a885b910"]
    // let orderKeys = ["0x98e25dd9a45bbf79100ebe3b1b311b2b6702a28c9fca5ee317feb0049893faa5",
    //     "0x0c78b81d5da49fe7fd13832aac4aba9f79f31d25453b61ed09ec3ce941adca70",
    //     "0x201dc11898ad0213485b4b34b9702beedc8f3bbcc71b2e38512508adb59c8ea9"];

    for (let i = 0; i < 1; i++) {
        let info = await getOrderInfo(orderKeys[i]);
        let sellOrder = info.order;
        // console.log("sellOrder: ", sellOrder);
        let buyOrder = {
            side: Side.Bid,
            saleKind: SaleKind.FixedPriceForItem,
            maker: trader.address,
            nft: sellOrder.nft,
            price: sellOrder.price,
            expiry: sellOrder.expiry,
            salt: sellOrder.salt,
        }

        let tx = await esDex.connect(trader).matchOrder(sellOrder, buyOrder, { value: toBn("0.002") });
        let txRec = await tx.wait();
        console.log("matchOrder tx: ", tx.hash);
    }

    // 5. else
    // await withdrawProtocolFee();
    // await testBatchTransferERC721();
}

/**
 * deployer:  0xBb03D7e8380a60665710582E42668cFceF85f457
trader:  0x36c4204DEEe1d4D44dAaB03f9052DB8F147b121F
Approval tx: 0x86adcab5ad5425ee06a4f4dfde541dd9957c45700ae40f2a339f187e1718e720
0x188fc9a7ed5f707d5014d34e8107d8e27ed413ca23b36a4c57585f7081cb5183
0x754b24b562ce0c9c4bbfdc33d32df978ea1aa24a05b5db5bc4276f2b117f40a7
0x502dc175e217a1674ea97c4879955f349601f1042fece72a6287e9d6d3e409d3
0x8cbf16bcb1b81cd09afc0251adacd640414e71c8de19da8838f2fab6873e570e
0x09dda3a84ece18ad761b938175e90c58e47db5afc6b5932128ea42a2ef9b029a
0x6656bc021b1103dace0acb064db8daa7a5826aa6855df0a082c6bf4a796c9f37
0x79d7f95e7a00b5e98f8148b658864a6fca3e9047dd02d9404fde23527976333e
0x5d3592bfb400e2d527bbbcad867815362e463889e9fc0bc7ae8e896109fbf008
0x6873f7adb43edc8ba7b1df84dc27d5cdbbe7682e76c8896d53eca2461f90f988
0xdc251c7d6e9d57ebfd3a72473d779cfb6888a66b2cd948adb19b8e4fa3cdfe6f
0xd3ca9e919cbc76db354eebf3fb8b0e6b44abf7a50aa62b85dd28f6ecdb7825fb
0x1b31abd2729f0a3954a54626ecf438dd0ca88fd421e952093fe8f28ca8da0cf2
0x867c70aae51725e71be494bc116f4c1116db6578c3bc2c28f20369bdfe587f95
0xce48e117d512125c4dd232f1f214018bbae5b81bc448ff55093cb11e5770abf8
0xa5b53aac6ebada4e71b6d44d7265d9b5d7b512a5dd841ddf34735ff5b7fbee1f
0x696c1ee9d1bc41dd3d5f8adc4292c6aee6f3937a9e8f06afc4325501d4c37ecb
0x6fda7f22aebe55f48b053c976cd8b222de6dd690b5d6586acc7e14705b5ff7b1
0x918e6b817fd4073ff77fc35c5cd9ddec5f36ae83c45bfded406c46608abf253e
0x7de730e1d3513e1b6e42e49bbcea494cd53f61ad5611a364788e4d88b5a26dc4
0x61d49cbcaea7da62f729f3ca43bccb28aebf29cffab2aeed0757fe6ae4403b07

 * 0x0e4b2f1af50024e46296a989ef65a468ee03f321b4e8a07ccc6dbb9e624194a8
 * 0x552f61ce9b7d8270fd19da3ff8f3c68857b630d8ce662b74caf34388b7118454
 * 
 * 
 * 0xfb97fe6ae1440f1b410d9678d0d066744f7c4cb348ee7fb8160c2ac7099e81c5
 * 
 */

async function approvalForVault() {
    // check is approved
    let isApproved = await testERC721.isApprovedForAll(deployer.address, esVault_address);

    if (isApproved) {
        console.log("Already approved");
        return;
    }

    let tx = await testERC721.setApprovalForAll(esVault_address, true);
    await tx.wait();
    console.log("Approval tx:", tx.hash)
}

async function testMakeOrder(tokenId = 0) {
    let now = parseInt(new Date() / 1000) + 100000
    let salt = 1;
    let nftAddress = erc721_address;
    // let tokenId = 0;
    let order = {
        side: Side.List,
        saleKind: SaleKind.FixedPriceForItem,
        maker: deployer.address,
        nft: [tokenId, nftAddress, 1],
        price: toBn("0.002"),
        expiry: now,
        salt: salt,
    }

    tx = await esDex.makeOrders([order]);
    txRec = await tx.wait();
    console.log(tx.hash);
}

async function testCancelOrder(orderKeys) {
    tx = await esDex.cancelOrders(orderKeys);
    txRec = await tx.wait();
    console.log(txRec);
}

async function testMatchOrder() {
    let now = 1734937947;
    let salt = 1;
    let tokenId = 0;
    let nftAddress = erc721_address;

    let sellOrder = {
        side: Side.List,
        saleKind: SaleKind.FixedPriceForItem,
        maker: deployer.address,
        nft: [tokenId, nftAddress, 1],
        price: toBn("0.002"),
        expiry: now,
        salt: salt,
    }

    // tx = await esDex.makeOrders([sellOrder]);
    // txRec = await tx.wait();
    // console.log("sellOrder tx: ", tx.hash);

    // ====
    let buyOrder = {
        side: Side.Bid,
        saleKind: SaleKind.FixedPriceForCollection,
        maker: trader.address,
        nft: [tokenId, nftAddress, 1],
        price: toBn("0.002"),
        expiry: now,
        salt: salt,
    }

    tx = await esDex.connect(trader).matchOrder(sellOrder, buyOrder, { value: toBn("0.002") });
    txRec = await tx.wait();
    console.log("matchOrder tx: ", txRec.hash);
}

async function testBatchTransferERC721() {
    toAddr = "0x7752A564c941f7145AdF8B50AA2eC975cEf58689"
    nftAddr = "0x3c8ac104dcbf03ae12c9ac80aa830e1b39609e97"
    tokenId = 1159
    asset = [nftAddr, tokenId]
    assets = [asset]
    tx = await esVault.callStatic.batchTransferERC721(toAddr, assets);
    console.log("tx: ", tx);
}

async function getOrderInfo(orderKey) {
    orderInfo = await esDex.orders(orderKey);
    // console.log("orderInfo: ", orderInfo);
    return orderInfo;
}

async function getfillsStat(orderKey) {
    fillStat = await esDex.filledAmount(orderKey);
    // console.log(fillStat);
    return fillStat;
}

async function withdrawProtocolFee() {
    await esDex.withdrawETH(deployer.address, toBn("0.00011"), { gasLimit: 100000 });
    console.log("WithdrawETH succeed.");

}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })
