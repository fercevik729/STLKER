const express = require('express');
const cors = require('cors');
const app = express();

app.use(cors());
app.use(express.json())

app.use('/login', (req, res) => {
  console.log(JSON.stringify(req.body))
  res.send({
    token: 'test123'
  });
});

app.listen(8080, () => console.log('API is running on http://localhost:8080/login'))
