const express = require("express");
const app = express();
const cors = require("cors");
const fetch = require("node-fetch");
require("dotenv").config();
const port = 3001;

app.use(cors());
app.use(express.json());

app.get("/tokenPrice", async (req, res) => {
  try {
    const { query } = req;
    
    // Fetch price for first token
    const optionsOne = {
      method: 'POST',
      headers: {
        accept: 'application/json',
        'content-type': 'application/json'
      },
      body: JSON.stringify({
        addresses: [
          {
            network: 'eth-mainnet',
            address: query.addressOne
          }
        ]
      })
    };

    // Fetch price for second token
    const optionsTwo = {
      method: 'POST',
      headers: {
        accept: 'application/json',
        'content-type': 'application/json'
      },
      body: JSON.stringify({
        addresses: [
          {
            network: 'eth-mainnet',
            address: query.addressTwo
          }
        ]
      })
    };

    // Make parallel requests for both token prices
    const [responseOne, responseTwo] = await Promise.all([
      fetch(`https://api.g.alchemy.com/prices/v1/${process.env.ALCHEMY_API_KEY}/tokens/by-address`, optionsOne)
        .then(res => res.json()),
      fetch(`https://api.g.alchemy.com/prices/v1/${process.env.ALCHEMY_API_KEY}/tokens/by-address`, optionsTwo)
        .then(res => res.json())
    ]);
    // Log both responses

    const usdPrices = {
      tokenOne: responseOne.data[0].prices[0].value,
      tokenTwo: responseTwo.data[0].prices[0].value,
      ratio: responseOne.data[0].prices[0].value / responseTwo.data[0].prices[0].value
      };

    return res.status(200).json(usdPrices);
  } catch (error) {
    console.error('Error fetching token prices:', error);
    return res.status(500).json({ error: 'Failed to fetch token prices' });
  }
});

app.listen(port, () => {
  console.log(`Listening for API Calls on port ${port}`);
});