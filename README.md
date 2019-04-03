[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/powered-by-netflix.svg)](https://forthebadge.com)

![Backpulse](https://files.backpulse.io/backpulse.png#cache2 "Backpulse.io")

# Backpulse core
Backpulse is an API Based / Headless CMS.  
Your site's content is accessible directly via our RESTful API, on any web framework and any device.  

## Installation
With a correctly configured Go toolchain:
```bash
go get github.com/backpulse/core
```

## Usage
First, you need to create a config.json using the `config.json.template` file.
* **URI** : MongoDB server address (_mongodb://..._)
* **Database** : MongoDB database name
* **Secret** : A secret key to encrypt JWT
* **GmailAddress** : A gmail address if you wish to send confirmation emails
* **GmailPassword** : The password associated with the gmail address obviously
* **StripeKey** : Your Stripe Key if you wish to integrate Stripe
* **BucketName** : Your Google Cloud Storage Bucket's name to store user files (images, binaries, plain text...)

You can also pass all these variables as environment variables:
* MONGODB_URI
* DATABASE
* SECRET
* GMAIL_ADDRESS
* GMAIL_PASSWORD
* STRIPE_KEY
* BUCKET_NAME  


**Note**: If a `config.json` file is found, it will override environment variable.

Then, you need to get your Google Service Account Key:
* Go to this [page](https://console.cloud.google.com/apis/credentials/serviceaccountkey).
* Create a new account with the Project -> Owner role.
* Download your private key as JSON.
* Move it to the root of this project.
* Rename it `google_credentials.json`.

You can also pass the content of this json file as an environment variable:

GOOGLE_APPLICATION_CREDENTIALS = `{"type": "service_account", "project_id": "projectID", ...}`

You're all set to run **Backpulse**!
```bash
go build -o backpulse && backpulse
```

**Note**: By default Backpulse runs on port 8000, but can be overridden with the `PORT` environment variable.

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.


## License 
[MIT](https://github.com/backpulse/core/blob/master/LICENSE)