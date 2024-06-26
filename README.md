# JubiBot Response Lambda

[![Go Lambda CI/CD](https://github.com/kenrms/jubibot-response-lambda/actions/workflows/main.yml/badge.svg?branch=main)](https://github.com/kenrms/jubibot-response-lambda/actions/workflows/main.yml)

This is a serverless AWS Lambda function that integrates with the OpenAI API to provide AI-generated responses to Discord messages.

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Introduction

This is a powerful tool that allows you to integrate AI-powered responses into your applications. By leveraging the OpenAI API, this service can generate contextual and relevant responses to user messages, enhancing the user experience and providing intelligent conversational capabilities.

## Features

- Seamless integration with the OpenAI API
- Serverless architecture using AWS Lambda
- Scalable and highly available
- Customizable response generation

## Installation

To use the AI Response Service, you'll need to deploy the AWS Lambda function to your AWS account. You can do this by following these steps:

1. Clone the repository: `git clone https://github.com/kenrms/jubibot-response-lambda.git`
2. Navigate to the project directory: `cd jubibot-response-lambda`
3. Deploy the Lambda function using the AWS CLI or your preferred deployment method.

Make sure to set the necessary environment variables, such as the OpenAI API key, before deploying the function.

## Usage

To use the AI Response Service, you'll need to send a POST request to the Lambda function's API Gateway endpoint. The request body should contain the message data that you want the AI to process.

Here's an example of how you can use the service:

```
POST /jubibot-response
{
    "channelId": "123422324",
    "messageContent": "Is that true?",
    "referencedMessageId": "522325233",
    "referencedMessageContent": "The sky is red",
    "referencedMessageAuthor": "janesmith112"
}
```

The function will then send the message data to the OpenAI API and return the generated response in the response body.

## Contributing

If you'd like to contribute to the AI Response Service, please follow these steps:

1. Fork the repository
2. Create a new branch for your feature or bug fix
3. Make your changes and commit them
4. Push your changes to your fork
5. Submit a pull request

## License

This project is licensed under the [MIT License](LICENSE).
