const { spawn } = require('child_process');

// Test the MCP server with a tool call
const server = spawn('./elastic-integration-docs-mcp', [], {
  stdio: ['pipe', 'pipe', 'inherit']
});

// Send initialize request
const initRequest = {
  jsonrpc: '2.0',
  id: 1,
  method: 'initialize',
  params: {
    protocolVersion: '2024-11-05',
    capabilities: {
      roots: {
        listChanged: true
      }
    },
    clientInfo: {
      name: 'test-client',
      version: '1.0.0'
    }
  }
};

console.log('Sending initialize request...');
server.stdin.write(JSON.stringify(initRequest) + '\n');

// Send tool call request
const toolRequest = {
  jsonrpc: '2.0',
  id: 2,
  method: 'tools/call',
  params: {
    name: 'get_service_info',
    arguments: {
      serviceName: 'nginx'
    }
  }
};

setTimeout(() => {
  console.log('Sending get_service_info tool call...');
  server.stdin.write(JSON.stringify(toolRequest) + '\n');
}, 100);

// Handle responses
server.stdout.on('data', (data) => {
  const lines = data.toString().trim().split('\n');
  lines.forEach(line => {
    if (line.trim()) {
      try {
        const response = JSON.parse(line);
        if (response.id === 2) {
          console.log('Tool response:');
          console.log(JSON.stringify(response, null, 2));
          if (response.result && response.result.content) {
            console.log('\nContent:');
            response.result.content.forEach(content => {
              console.log(content.text);
            });
          }
        }
      } catch (e) {
        console.log('Raw output:', line);
      }
    }
  });
});

// Clean up after 3 seconds
setTimeout(() => {
  server.kill();
  process.exit(0);
}, 3000);
