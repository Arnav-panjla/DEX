const options = {
  method: 'POST',
  headers: {accept: 'application/json', 'content-type': 'application/json'},
  body: JSON.stringify({
    addresses: [
      {network: 'eth-mainnet', address: '0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48'}
    ]
  })
};

fetch('https://api.g.alchemy.com/prices/v1/FufLRoDAv02_l2dceqZxkiDC12Mo4MIi/tokens/by-address', options)
    .then(res => res.json())  // Parse JSON response
    .then(data => {
    // Access the prices array in the response
        console.log(data.data[0].prices);
    })
    .catch(err => console.error('Error:', err)); 