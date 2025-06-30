import openai

client = openai.OpenAI(
    base_url="https://api.llm7.io/v1",
    api_key="unused"  # Or get it for free at https://token.llm7.io/ for higher rate limits.
)

response = client.chat.completions.create(
    model="gpt-4.1-nano",
    messages=[
        {"role": "user", "content": "Tell me a short story about a brave squirrel."}
    ]
)

print(response.choices[0].message.content)
