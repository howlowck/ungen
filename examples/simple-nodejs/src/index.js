// simple hello world app
const express = require('express');
const app = express();
// UNGEN: copy next 1 line to cb.port
const port = 3000;

// UNGEN: replace "World" with substitute(upperCase(kebabCase(var.appName)), "-", ".")
app.get('/', (req, res) => res.send('Hello World!'));

// UNGEN: insert substitute(cb.port, "port", "newPort")

// start the Express server
app.listen(port, () => {
    console.log(`server started at http://localhost:${port}`);
});
