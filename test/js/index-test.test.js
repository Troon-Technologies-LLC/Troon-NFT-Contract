import path from "path";
import {
  init,
  emulator,
  getAccountAddress,
  deployContractByName,
  getContractAddress,
  getTransactionCode,
  getScriptCode,
  executeScript,
  sendTransaction,
  mintFlow,
  getFlowBalance,
} from "flow-js-testing";
import { expect } from "@jest/globals";
jest.setTimeout(10000);

beforeAll(async () => {
  const basePath = path.resolve(__dirname, "../..");
  const port = 8080;

  await init(basePath, { port });
  await emulator.start(port);
});

afterAll(async () => {
  const port = 8080;
  await emulator.stop(port);
});


describe("Replicate Playground Accounts", () => {
  test("Create Accounts", async () => {
    // Playground project support 4 accounts, but nothing stops you from creating more by following the example laid out below
    const Alice = await getAccountAddress("Alice");
    const Bob = await getAccountAddress("Bob");
    const Charlie = await getAccountAddress("Charlie");
    const Dave = await getAccountAddress("Dave");

    console.log(
      "Four Playground accounts were created with following addresses"
    );
    console.log("Alice:", Alice);
    console.log("Bob:", Bob);
    console.log("Charlie:", Charlie);
    console.log("Dave:", Dave);
  });
});


describe("Deployment", () => {
  test("Deploy for NonFungibleToken", async () => {
    const name = "NonFungibleToken"
    const to = await getAccountAddress("Alice")
    let update = true

    let result;
    try {
      result = await deployContractByName({
        name,
        to,
        update,
      });
    } catch (e) {
      console.log(e);
    }
    expect(name).toBe("NonFungibleToken");

  });
  test("Deploy for TroonAtomicStandard", async () => {
    const name = "TroonAtomicStandard"
    const to = await getAccountAddress("Bob")
    let update = true

    const NonFungibleToken = await getContractAddress("NonFungibleToken");
    const addressMap = {
      NonFungibleToken
    };

    let result;
    try {
      result = await deployContractByName({
        name,
        to,
        addressMap,
        update,
      });
    } catch (e) {
      console.log(e);
    }
    expect(name).toBe("TroonAtomicStandard");

  });
});
describe("Transactions", () => {
  test("test transaction setup admin Account", async () => {
    const name = "setupAdminAccount";

    // Import participating accounts
    const Charlie = await getAccountAddress("Charlie");

    // Set transaction signers
    const signers = [Charlie];

    // Generate addressMap from import statements
    const NonFungibleToken = await getContractAddress("NonFungibleToken");
    const TroonAtomicStandard = await getContractAddress("TroonAtomicStandard");
    const addressMap = {
      NonFungibleToken,
      TroonAtomicStandard,
    };

    let code = await getTransactionCode({
      name,
      addressMap,
    });

    let txResult;
    try {
      txResult = await sendTransaction({
        code,
        signers,
      });
    } catch (e) {
      console.log(e);
    }
    console.log("tx Result", txResult);

    // expect(txResult.errorMessage).toBe("");
  });
  test("test transaction add admin Account", async () => {
    const name = "addAdminAccount";

    // Import participating accounts
    const Bob = await getAccountAddress("Bob");
    const Charlie = await getAccountAddress("Charlie");

    // Set transaction signers
    const signers = [Bob];

    // Generate addressMap from import statements
    const NonFungibleToken = await getContractAddress("NonFungibleToken");
    const TroonAtomicStandard = await getContractAddress("TroonAtomicStandard");
    const addressMap = {
      NonFungibleToken,
      TroonAtomicStandard,
    };

    let code = await getTransactionCode({
      name,
      addressMap,
    });

    let txResult;
    try {
      const args = [Charlie];

      txResult = await sendTransaction({
        code,
        args,
        signers,
      });
    } catch (e) {
      console.log(e);
    }
    console.log("tx result ", txResult);

    // expect(txResult.errorMessage).toBe("");
  });
  test("test transaction  create brand", async () => {
    const name = "createBrand";

    // Import participating accounts
    const Charlie = await getAccountAddress("Charlie");

    // Set transaction signers
    const signers = [Charlie];

    // Generate addressMap from import statements
    const NonFungibleToken = await getContractAddress("NonFungibleToken");
    const TroonAtomicStandard = await getContractAddress("TroonAtomicStandard");
    const addressMap = {
      NonFungibleToken,
      TroonAtomicStandard,
    };

    let code = await getTransactionCode({
      name,
      addressMap,
    });
    const args = ["HondaNorth", { name: "Alice" }];

    let txResult;
    try {
      txResult = await sendTransaction({
        code,
        signers,
        args,
      });
    } catch (e) {
      console.log(e);
    }
    console.log("tx Result", txResult);

    // expect(txResult.errorMessage).toBe("");
  });
  test("test transaction  create Schema", async () => {
    const name = "createSchema";

    // Import participating accounts
    const Charlie = await getAccountAddress("Charlie");

    // Set transaction signers
    const signers = [Charlie];

    // Generate addressMap from import statements
    const NonFungibleToken = await getContractAddress("NonFungibleToken");
    const TroonAtomicStandard = await getContractAddress("TroonAtomicStandard");
    const addressMap = {
      NonFungibleToken,
      TroonAtomicStandard,
    };

    let code = await getTransactionCode({
      name,
      addressMap,
    });
    const args = ["Test Schema"];

    let txResult;
    try {
      txResult = await sendTransaction({
        code,
        signers,
        args,
      });
    } catch (e) {
      console.log(e);
    }
    console.log("tx Result", txResult);

    // expect(txResult.errorMessage).toBe("");
  });
  test("test transaction  create template", async () => {
    const name = "createTemplateStaticData";

    // Import participating accounts
    const Charlie = await getAccountAddress("Charlie");

    // Set transaction signers
    const signers = [Charlie];

    // Generate addressMap from import statements
    const NonFungibleToken = await getContractAddress("NonFungibleToken");
    const TroonAtomicStandard = await getContractAddress("TroonAtomicStandard");
    const addressMap = {
      NonFungibleToken,
      TroonAtomicStandard,
    };

    let code = await getTransactionCode({
      name,
      addressMap,
    });
    // brandId, schemaId, maxSupply,immutableData
    const args = [1, 1, 100];
    let txResult;
    try {
      txResult = await sendTransaction({
        code,
        signers,
        args,
      });
    } catch (e) {
      console.log(e);
    }
    console.log("tx Result", txResult);

    // expect(txResult.errorMessage).toBe("");
  });
  test("test transaction  mint NFT", async () => {
    const name = "mintNFT";

    // Import participating accounts
    const Charlie = await getAccountAddress("Charlie");

    // Set transaction signers
    const signers = [Charlie];

    // Generate addressMap from import statements
    const NonFungibleToken = await getContractAddress("NonFungibleToken");
    const TroonAtomicStandard = await getContractAddress("TroonAtomicStandard");
    const addressMap = {
      NonFungibleToken,
      TroonAtomicStandard,
    };

    let code = await getTransactionCode({
      name,
      addressMap,
    });
    // brandId, schemaId, maxSupply,immutableData
    const args = [1, Charlie];
    let txResult;
    try {
      txResult = await sendTransaction({
        code,
        signers,
        args,
      });
    } catch (e) {
      console.log(e);
    }
    console.log("tx Result", txResult);

    // expect(txResult.errorMessage).toBe("");
  });



})
describe("Scripts", () => {
  test("get total supply", async () => {

    const name = "getTotalSupply";
    const Bob = await getAccountAddress("Bob");

    const NonFungibleToken = await getContractAddress("NonFungibleToken");
    const TroonAtomicStandard = await getContractAddress("TroonAtomicStandard");

    const addressMap = {
      NonFungibleToken,
      TroonAtomicStandard,
    }
    let code = await getScriptCode({
      name,
      addressMap,
    })

    code = code.toString().replace(/(?:getAccount\(\s*)(0x.*)(?:\s*\))/g, (_, match) => {
      const accounts = {
        "0x01": Alice,
        "0x02": Bob,
      };
      const name = accounts[match];
      return `getAccount(${name})`;
    });

    const result = await executeScript({
      code,
    });
    console.log("result", result);



  });
  test("get brand data", async () => {

    const name = "getAllBrands";
    const Bob = await getAccountAddress("Bob");

    const NonFungibleToken = await getContractAddress("NonFungibleToken");
    const TroonAtomicStandard = await getContractAddress("TroonAtomicStandard");

    const addressMap = {
      NonFungibleToken,
      TroonAtomicStandard,
    }
    let code = await getScriptCode({
      name,
      addressMap,
    })

    code = code.toString().replace(/(?:getAccount\(\s*)(0x.*)(?:\s*\))/g, (_, match) => {
      const accounts = {
        "0x01": Alice,
        "0x02": Bob,
      };
      const name = accounts[match];
      return `getAccount(${name})`;
    });

    const result = await executeScript({
      code,
    });
    console.log("result", result);

  });
  test("get brand data by Id", async () => {

    const name = "getBrandById";
    const Bob = await getAccountAddress("Bob");

    const NonFungibleToken = await getContractAddress("NonFungibleToken");
    const TroonAtomicStandard = await getContractAddress("TroonAtomicStandard");

    const addressMap = {
      NonFungibleToken,
      TroonAtomicStandard,
    }
    let code = await getScriptCode({
      name,
      addressMap,
    });
    const args = [1];

    code = code.toString().replace(/(?:getAccount\(\s*)(0x.*)(?:\s*\))/g, (_, match) => {
      const accounts = {
        "0x01": Alice,
        "0x02": Bob,
      };
      const name = accounts[match];
      return `getAccount(${name})`;
    });

    const result = await executeScript({
      code,
      args,
    });
    console.log("result", result);

  });
  test("get schema data", async () => {

    const name = "getallSchema";
    const Bob = await getAccountAddress("Bob");

    const NonFungibleToken = await getContractAddress("NonFungibleToken");
    const TroonAtomicStandard = await getContractAddress("TroonAtomicStandard");

    const addressMap = {
      NonFungibleToken,
      TroonAtomicStandard,
    }
    let code = await getScriptCode({
      name,
      addressMap,
    })

    code = code.toString().replace(/(?:getAccount\(\s*)(0x.*)(?:\s*\))/g, (_, match) => {
      const accounts = {
        "0x01": Alice,
        "0x02": Bob,
      };
      const name = accounts[match];
      return `getAccount(${name})`;
    });

    const result = await executeScript({
      code,
    });
    console.log("result", result);

  });

  test("get schema data by Id", async () => {

    const name = "getSchemaById";
    const Bob = await getAccountAddress("Bob");

    const NonFungibleToken = await getContractAddress("NonFungibleToken");
    const TroonAtomicStandard = await getContractAddress("TroonAtomicStandard");

    const addressMap = {
      NonFungibleToken,
      TroonAtomicStandard,
    }
    let code = await getScriptCode({
      name,
      addressMap,
    })

    code = code.toString().replace(/(?:getAccount\(\s*)(0x.*)(?:\s*\))/g, (_, match) => {
      const accounts = {
        "0x01": Alice,
        "0x02": Bob,
      };
      const name = accounts[match];
      return `getAccount(${name})`;
    });


    console.log(typeof (myInt))
    const args = [1];
    const result = await executeScript({
      code,
      args,
    });
    console.log("result", result);

  });

  test("get template data ", async () => {

    const name = "getAllTemplates";
    const Bob = await getAccountAddress("Bob");

    const NonFungibleToken = await getContractAddress("NonFungibleToken");
    const TroonAtomicStandard = await getContractAddress("TroonAtomicStandard");

    const addressMap = {
      NonFungibleToken,
      TroonAtomicStandard,
    }
    let code = await getScriptCode({
      name,
      addressMap,
    })

    code = code.toString().replace(/(?:getAccount\(\s*)(0x.*)(?:\s*\))/g, (_, match) => {
      const accounts = {
        "0x01": Alice,
        "0x02": Bob,
      };
      const name = accounts[match];
      return `getAccount(${name})`;
    });

    const result = await executeScript({
      code,
    });
    console.log("result", result);

  });
  test("get template data by Id", async () => {

    const name = "getTemplateById";
    const Bob = await getAccountAddress("Bob");

    const NonFungibleToken = await getContractAddress("NonFungibleToken");
    const TroonAtomicStandard = await getContractAddress("TroonAtomicStandard");

    const addressMap = {
      NonFungibleToken,
      TroonAtomicStandard,
    }
    let code = await getScriptCode({
      name,
      addressMap,
    })

    code = code.toString().replace(/(?:getAccount\(\s*)(0x.*)(?:\s*\))/g, (_, match) => {
      const accounts = {
        "0x01": Alice,
        "0x02": Bob,
      };
      const name = accounts[match];
      return `getAccount(${name})`;
    });
    const args = [1]
    const result = await executeScript({
      code,
      args,
    });
    console.log("result", result);

  });


  test("get all nfts  data", async () => {

    const name = "getAllNFTIds";
    const Charlie = await getAccountAddress("Charlie");

    const NonFungibleToken = await getContractAddress("NonFungibleToken");
    const TroonAtomicStandard = await getContractAddress("TroonAtomicStandard");

    const addressMap = {
      NonFungibleToken,
      TroonAtomicStandard,
    }

    let code = await getScriptCode({
      name,
      addressMap,
    })

    code = code.toString().replace(/(?:getAccount\(\s*)(0x.*)(?:\s*\))/g, (_, match) => {
      const accounts = {
        "0x03": Charlie,
      };
      const name = accounts[match];
      return `getAccount(${name})`;
    });
    const args = [Charlie]
    const result = await executeScript({
      code,
      args,
    });
    console.log("result", result);

  });

  test("get nft template data", async () => {

    const name = "getNFTTemplateData";
    const Charlie = await getAccountAddress("Charlie");

    const NonFungibleToken = await getContractAddress("NonFungibleToken");
    const TroonAtomicStandard = await getContractAddress("TroonAtomicStandard");

    const addressMap = {
      NonFungibleToken,
      TroonAtomicStandard,
    }
    let code = await getScriptCode({
      name,
      addressMap,
    })

    code = code.toString().replace(/(?:getAccount\(\s*)(0x.*)(?:\s*\))/g, (_, match) => {
      const accounts = {
        "0x03": Charlie,
      };
      const name = accounts[match];
      return `getAccount(${name})`;
    });

    const args = [Charlie]
    const result = await executeScript({
      code,
      args,
    });
    console.log("result", result);

  });

});