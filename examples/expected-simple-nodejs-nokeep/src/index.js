// simple hello world app
const express = require('express');
const app = express();
const port = 3000;

app.get('/', (req, res) => res.send('Hello HAOS.AWESOME.APP!!'));

const newPort = 3000;

// start the Express server
app.listen(port, () => {
    console.log(`server started at http://localhost:${port}`);
});
