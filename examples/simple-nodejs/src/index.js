// simple hello world app
const express = require('express');
const app = express();
const port = 3000;

// UNGEN: replace "World" with kebabCase(var.appName)
app.get('/', (req, res) => res.send('Hello World!'));

// start the Express server
app.listen(port, () => {
    console.log(`server started at http://localhost:${port}`);
});
