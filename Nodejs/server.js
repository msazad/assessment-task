const http = require('http');
const fs = require('fs');
const path = require('path');

// Function to generate primary and secondary MAC addresses
function generateMAC(cID) {
    const last64 = cID.slice(-16);
    const num = BigInt(`0x${last64}`);
    const primaryMAC = (num & BigInt('0xFFFFFFFFFFFF')) | BigInt('0x020000000000'); // Set the unicast bit
    const mac1 = formatMAC(primaryMAC);
    const mac2 = formatMAC((primaryMAC + BigInt(1)) & BigInt('0xFFFFFFFFFFFF'));
    return { mac1, mac2 };
}

// Function to format the MAC address into a human-readable form
function formatMAC(num) {
    return [
        (num >> BigInt(40)) & BigInt(0xFF),
        (num >> BigInt(32)) & BigInt(0xFF),
        (num >> BigInt(24)) & BigInt(0xFF),
        (num >> BigInt(16)) & BigInt(0xFF),
        (num >> BigInt(8)) & BigInt(0xFF),
        num & BigInt(0xFF)
    ].map(byte => byte.toString(16).padStart(2, '0')).join(':');
}

// Create the server
http.createServer((req, res) => {
    if (req.method === 'POST' && req.url === '/generate') {
        let body = '';

        req.on('data', chunk => {
            body += chunk.toString();
        });

        req.on('end', () => {
            const { cid } = JSON.parse(body);
            if (!/^[0-9a-fA-F]{32}$/.test(cid)) {
                res.writeHead(400, { 'Content-Type': 'application/json' });
                return res.end(JSON.stringify({ error: 'Invalid 128-bit hexadecimal CID. Ensure it is 32 characters long.' }));
            }

            const { mac1, mac2 } = generateMAC(cid);
            res.writeHead(200, { 'Content-Type': 'application/json' });
            res.end(JSON.stringify({ mac1, mac2 }));
        });
    } else if (req.method === 'GET' && req.url === '/') {
        // Serve the index.html file
        const filePath = path.join(__dirname, 'index.html');
        fs.readFile(filePath, (err, data) => {
            if (err) {
                res.writeHead(500);
                return res.end('Error loading index.html');
            }
            res.writeHead(200, { 'Content-Type': 'text/html' });
            res.end(data);
        });
    } else {
        res.writeHead(404);
        res.end('404 Not Found');
    }
}).listen(8081, () => {
    console.log('Node.js server running on http://localhost:8081');
});
