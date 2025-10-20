import asyncio
import os
import json
from contextlib import AsyncExitStack
from typing import Optional, List, Dict, Any

from mcp import ClientSession, StdioServerParameters, Message
from mcp.client.stdio import stdio_client
from openai import AsyncOpenAI
from dotenv import load_dotenv

# --- Configuration ---
# Load environment variables from a .env file (e.g., OPENAI_API_KEY)
load_dotenv()

# IMPORTANT: Replace this with the command you use to run your MCP server.
# Assuming your server is a Python script named 'elastic_server.py'
# and the server's main function is decorated to expose tools.
# If your server is an executable, just use the executable name.
SERVER_COMMAND = ["python", "elastic_server.py"]
SERVER_NAME = "elastic-integration-docs-mcp" # Must match the name defined in your server
MODEL_NAME = "gpt-4o" # Ensure your chosen model supports tool calling

class MCPClientAgent:
    """
    An interactive client agent that connects to an MCP server via stdio
    and uses an LLM (OpenAI) to decide which tools to call.
    """
    def __init__(self, command: List[str], server_name: str, model_name: str):
        self.server_command = command
        self.server_name = server_name
        self.model_name = model_name
        self.client: Optional[AsyncOpenAI] = None
        self.session: Optional[ClientSession] = None
        self.exit_stack = AsyncExitStack()
        self.history: List[Dict[str, Any]] = []

    async def initialize(self):
        """Initializes the OpenAI client and connects to the MCP server."""
        print("Initializing MCP Client...")
        
        # 1. Initialize OpenAI Client
        self.client = AsyncOpenAI(api_key=os.getenv("OPENAI_API_KEY"))

        # 2. Set up the stdio server parameters
        # This tells the client how to launch and communicate with the server process.
        server_params = StdioServerParameters(command=self.server_command)

        # 3. Connect to the server
        self.session = await self.exit_stack.enter_async_context(
            stdio_client(
                server_parameters=server_params,
                server_name=self.server_name,
                timeout=30 # Give up to 30 seconds for the server to launch and register
            )
        )
        print(f"Connection established to MCP server: {self.server_name}")
        
        # 4. Get tool definitions from the server
        # The session.tools() returns the JSON schema for the tools.
        self.tool_definitions = [tool.json_schema for tool in self.session.tools()]
        print(f"Found {len(self.tool_definitions)} tool(s) on the server.")
        
        # 5. Add a system instruction to guide the LLM
        system_prompt = (
            "You are an expert AI agent designed to answer user questions about Elastic Integration "
            "documentation. Use the provided tools exclusively to fetch information. "
            "Do not guess. Only use the tools when necessary. If a tool call is needed, "
            "only output the tool call and do not add any conversational text."
        )
        self.history.append({"role": "system", "content": system_prompt})

        print("\nReady for interaction. Type 'exit' to quit.")

    async def process_user_query(self, query: str):
        """Handles the main interaction loop with the LLM and tool execution."""

        # 1. Add user query to history
        self.history.append({"role": "user", "content": query})

        try:
            while True:
                # 2. Call the LLM with the current conversation history and tools
                completion = await self.client.chat.completions.create(
                    model=self.model_name,
                    messages=self.history,
                    tools=self.tool_definitions,
                    tool_choice="auto", # Allow the model to decide whether to call a tool
                )

                response_message = completion.choices[0].message

                # 3. Check for tool calls
                if response_message.tool_calls:
                    print("\n[AGENT] 🤖 Tool call requested...")

                    tool_calls_to_append = []

                    # 4. Execute all requested tool calls
                    for tool_call in response_message.tool_calls:
                        function_name = tool_call.function.name
                        function_args = json.loads(tool_call.function.arguments)

                        print(f"  -> Calling tool '{function_name}' with args: {function_args}")

                        # Execute the tool call using the MCP client session
                        # self.session.call_tool() handles sending the request to the stdio server
                        tool_output = await self.session.call_tool(
                            function_name,
                            **function_args
                        )

                        # The tool output is wrapped in an MCP Message, extract the data
                        if isinstance(tool_output, Message):
                            tool_result_content = tool_output.data
                        else:
                            tool_result_content = tool_output

                        print(f"  -> Tool result received (Type: {type(tool_result_content).__name__})")

                        # 5. Append the LLM's request and the tool's output to the history
                        # The client must send the tool call response back to the LLM
                        tool_calls_to_append.append({
                            "tool_call_id": tool_call.id,
                            "function": {"name": function_name},
                            "content": json.dumps(tool_result_content)
                        })

                    # Add the LLM's request to the history
                    self.history.append(response_message)
                    # Add the tool results to the history
                    self.history.append({
                        "role": "tool",
                        "tool_calls": tool_calls_to_append # This is the array of results
                    })

                    # Loop again to send the tool results back to the LLM
                    # and get the final natural language response
                    continue

                # 6. If no tool calls, it's the final answer
                final_answer = response_message.content
                print(f"\n[AGENT] ✅ Final Answer:\n{final_answer}")

                # Update history with the final response
                self.history.append(response_message)
                break # Exit the loop after the final answer

        except Exception as e:
            print(f"\n[ERROR] An error occurred: {e}")
            # Remove the last user message from history on failure
            self.history.pop()

    async def run_interactive_session(self):
        """Runs the main interactive terminal session."""
        try:
            await self.initialize()
            while True:
                user_input = input("\n[YOU] > ")
                if user_input.lower() in ['exit', 'quit']:
                    print("Exiting interactive session.")
                    break

                if user_input.strip():
                    await self.process_user_query(user_input.strip())
        except Exception as e:
            print(f"\n[CRITICAL ERROR] Failed to run client: {e}")
        finally:
            await self.exit_stack.aclose()


if __name__ == "__main__":
    # Ensure you have your MCP server file ready (e.g., elastic_server.py)
    # and a .env file with your OPENAI_API_KEY.
    agent = MCPClientAgent(
        command=SERVER_COMMAND,
        server_name=SERVER_NAME,
        model_name=MODEL_NAME
    )

    # Run the asynchronous function
    asyncio.run(agent.run_interactive_session())

