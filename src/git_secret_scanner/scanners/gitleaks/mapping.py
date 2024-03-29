from git_secret_scanner.report import SecretKind


GITLEAKS_RULE_TO_SECRET_KIND = {
    'adafruit-api-key': SecretKind.AdafruitIO,
    'adobe-client-id': SecretKind.AdobeIO,
    'adobe-client-secret': SecretKind.AdobeIO,
    'airtable-api-key': SecretKind.Airtable,
    'algolia-api-key': SecretKind.Algolia,
    'alibaba-access-key-id': SecretKind.Alibaba,
    'alibaba-secret-key': SecretKind.Alibaba,
    'asana-client-id': SecretKind.Asana,
    'asana-client-secret': SecretKind.Asana,
    'atlassian-api-token': SecretKind.Atlassian,
    'authress-service-client-access-key': SecretKind.Authress,
    'aws-access-token': SecretKind.AWS,
    'beamer-api-token': SecretKind.Beamer,
    'bitbucket-client-id': SecretKind.BitBucket,
    'bitbucket-client-secret': SecretKind.BitBucket,
    'bittrex-access-key': SecretKind.Bittrex,
    'bittrex-secret-key': SecretKind.Bittrex,
    'clojars-api-token': SecretKind.Clojars,
    'codecov-access-token': SecretKind.Codecov,
    'coinbase-access-token': SecretKind.Coinbase,
    'confluent-access-token': SecretKind.Confluent,
    'confluent-secret-key': SecretKind.Confluent,
    'contentful-delivery-api-token': SecretKind.Confluent,
    'databricks-api-token': SecretKind.Databricks,
    'datadog-access-token': SecretKind.Datadog,
    'defined-networking-api-token': SecretKind.DefinedNetworking,
    'digitalocean-access-token': SecretKind.DigitalOcean,
    'digitalocean-pat': SecretKind.DigitalOcean,
    'digitalocean-refresh-token': SecretKind.DigitalOcean,
    'discord-api-token': SecretKind.Discord,
    'discord-client-id': SecretKind.Discord,
    'discord-client-secret': SecretKind.Discord,
    'doppler-api-token': SecretKind.Doppler,
    'droneci-access-token': SecretKind.DroneCI,
    'dropbox-api-token': SecretKind.Dropbox,
    'dropbox-long-lived-api-token': SecretKind.Dropbox,
    'dropbox-short-lived-api-token': SecretKind.Dropbox,
    'duffel-api-token': SecretKind.Duffel,
    'dynatrace-api-token': SecretKind.Dynatrace,
    'easypost-api-token': SecretKind.Generic,
    'easypost-test-api-token': SecretKind.Generic,
    'etsy-access-token': SecretKind.Etsy,
    'facebook': SecretKind.Facebook,
    'fastly-api-token': SecretKind.Fastly,
    'finicity-api-token': SecretKind.Finicity,
    'finicity-client-secret': SecretKind.Finicity,
    'finnhub-access-token': SecretKind.Finnhub,
    'flickr-access-token': SecretKind.Flickr,
    'flutterwave-encryption-key': SecretKind.Flutterwave,
    'flutterwave-public-key': SecretKind.Flutterwave,
    'flutterwave-secret-key': SecretKind.Flutterwave,
    'frameio-api-token': SecretKind.FrameIO,
    'freshbooks-access-token': SecretKind.Freshbooks,
    'gcp-api-key': SecretKind.GCP,
    'generic-api-key': SecretKind.Generic,
    'github-app-token': SecretKind.Github,
    'github-fine-grained-pat': SecretKind.Github,
    'github-oauth': SecretKind.Github,
    'github-pat': SecretKind.Github,
    'github-refresh-token': SecretKind.Github,
    'gitlab-pat': SecretKind.Gitlab,
    'gitlab-ptt': SecretKind.Gitlab,
    'gitlab-rrt': SecretKind.Gitlab,
    'gitter-access-token': SecretKind.Gitter,
    'gocardless-api-token': SecretKind.GoCardless,
    'grafana-api-key': SecretKind.Grafana,
    'grafana-cloud-api-token': SecretKind.Grafana,
    'grafana-service-account-token': SecretKind.Grafana,
    'hashicorp-tf-api-token': SecretKind.TerraformCloud,
    'heroku-api-key': SecretKind.Heroku,
    'hubspot-api-key': SecretKind.HubSpot,
    'huggingface-access-token': SecretKind.HuggingFace,
    'huggingface-organization-api-token': SecretKind.HuggingFace,
    'intercom-api-key': SecretKind.Intercom,
    'jfrog-api-key': SecretKind.JFrog,
    'jfrog-identity-token':  SecretKind.JFrog,
    'jwt': SecretKind.JWT,
    'jwt-base64': SecretKind.JWT,
    'kraken-access-token': SecretKind.Kraken,
    'kucoin-access-token': SecretKind.KuCoin,
    'kucoin-secret-key': SecretKind.KuCoin,
    'launchdarkly-access-token': SecretKind.LaunchDarkly,
    'linear-api-key': SecretKind.Linear,
    'linear-client-secret': SecretKind.Linear,
    'linkedin-client-id': SecretKind.LinkedIn,
    'linkedin-client-secret': SecretKind.LinkedIn,
    'lob-api-key': SecretKind.Lob,
    'lob-pub-api-key': SecretKind.Lob,
    'mailchimp-api-key': SecretKind.Mailchimp,
    'mailgun-private-api-token': SecretKind.Mailgun,
    'mailgun-pub-key': SecretKind.Mailgun,
    'mailgun-signing-key': SecretKind.Mailgun,
    'mapbox-api-token': SecretKind.MapBox,
    'mattermost-access-token': SecretKind.Mattermost,
    'messagebird-api-token': SecretKind.MessageBird,
    'messagebird-client-id': SecretKind.MessageBird,
    'microsoft-teams-webhook': SecretKind.MicrosoftTeams,
    'netlify-access-token': SecretKind.Netlify,
    'new-relic-browser-api-token': SecretKind.NewRelic,
    'new-relic-user-api-id': SecretKind.NewRelic,
    'new-relic-user-api-key': SecretKind.NewRelic,
    'npm-access-token': SecretKind.Npm,
    'nytimes-access-token': SecretKind.Nytimes,
    'okta-access-token': SecretKind.Okta,
    'openai-api-key': SecretKind.OpenAI,
    'plaid-api-token': SecretKind.Plaid,
    'plaid-client-id': SecretKind.Plaid,
    'plaid-secret-key': SecretKind.Plaid,
    'planetscale-api-token': SecretKind.PlanetScale,
    'planetscale-oauth-token': SecretKind.PlanetScale,
    'planetscale-password': SecretKind.PlanetScale,
    'postman-api-token': SecretKind.Postman,
    'prefect-api-token': SecretKind.Prefect,
    'private-key': SecretKind.PrivateKey,
    'pulumi-api-token': SecretKind.Pulumi,
    'pypi-upload-token': SecretKind.PyPI,
    'rapidapi-access-token': SecretKind.RapidApi,
    'readme-api-token': SecretKind.ReadMe,
    'rubygems-api-token': SecretKind.RubyGems,
    'scalingo-api-token': SecretKind.Scalingo,
    'sendbird-access-id': SecretKind.Sendbird,
    'sendbird-access-token': SecretKind.Sendbird,
    'sendgrid-api-token': SecretKind.SendGrid,
    'sendinblue-api-token': SecretKind.SendinBlue,
    'sentry-access-token': SecretKind.Sentry,
    'shippo-api-token': SecretKind.Shippo,
    'shopify-access-token': SecretKind.Shopify,
    'shopify-custom-access-token': SecretKind.Shopify,
    'shopify-private-app-access-token': SecretKind.Shopify,
    'shopify-shared-secret': SecretKind.Shopify,
    'sidekiq-secret': SecretKind.Sidekiq,
    'sidekiq-sensitive-url': SecretKind.Sidekiq,
    'slack-app-token': SecretKind.Slack,
    'slack-bot-token': SecretKind.Slack,
    'slack-config-access-token': SecretKind.Slack,
    'slack-config-refresh-token': SecretKind.Slack,
    'slack-legacy-bot-token': SecretKind.Slack,
    'slack-legacy-token': SecretKind.Slack,
    'slack-legacy-workspace-token': SecretKind.Slack,
    'slack-user-token': SecretKind.Slack,
    'slack-webhook-url': SecretKind.Slack,
    'snyk-api-token': SecretKind.Snyk,
    'square-access-token': SecretKind.Square,
    'squarespace-access-token': SecretKind.Squarespace,
    'stripe-access-token': SecretKind.Stripe,
    'sumologic-access-id': SecretKind.SumoLogic,
    'sumologic-access-token': SecretKind.SumoLogic,
    'telegram-bot-api-token': SecretKind.Telegram,
    'travisci-access-token': SecretKind.TravisCI,
    'twilio-api-key': SecretKind.Twilio,
    'twitch-api-token': SecretKind.Twitch,
    'twitter-access-secret': SecretKind.Twitter,
    'twitter-access-token': SecretKind.Twitter,
    'twitter-api-key': SecretKind.Twitter,
    'twitter-api-secret': SecretKind.Twitter,
    'twitter-bearer-token': SecretKind.Twitter,
    'typeform-api-token': SecretKind.Typeform,
    'vault-batch-token': SecretKind.Vault,
    'vault-service-token': SecretKind.Vault,
    'yandex-access-token': SecretKind.Yandex,
    'yandex-api-key': SecretKind.Yandex,
    'yandex-aws-access-token': SecretKind.Yandex,
    'zendesk-secret-key': SecretKind.Zendesk,
}
