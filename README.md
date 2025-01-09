# Vultr Certbot DNS-01 Challenge Plugin

A lightweight and efficient plugin to automate DNS-01 challenges for **Let's Encrypt** SSL certificates using Vultr's DNS API.

## **Prerequisites**
Before getting started, ensure the following requirements are met:

1. **Vultr API Key**:
    - Obtain your API key from [Vultr's API settings](https://my.vultr.com/settings/#settingsapi).
2. **IP Whitelisting**:
    - Your machine's IP address must be whitelisted in the Vultr API settings to make API requests.

---

## **Installation**

1. Download the prebuilt binary for your platform:
   ```bash
   curl -sL https://github.com/h4ckitt/vultr-acme-dns/releases/download/v1.0.0/vultr-dns-linux-amd64 -o vultr-dns
   ```
2. Make the binary executable:
   ```bash
   chmod +x vultr-dns
   ```

---

## **Usage Steps**

### Step 1: Set Your Vultr API Key
Save your Vultr API key as an environment variable in a `.certbotenv` file.  
Example:
```bash
echo 'export VULTR_API_KEY=your_api_key_here' > .certbotenv
```

### Step 2: Obtain a Wildcard SSL Certificate
Run the following Certbot command to generate a wildcard certificate:
```bash
sudo certbot certonly --manual --preferred-challenges dns   --manual-auth-hook ./vultr-dns   --manual-cleanup-hook ./vultr-dns   -d "*.yourdomain.tld" -d "yourdomain.tld"
```

Replace:
- `*.yourdomain.tld`: Your wildcard domain.
- `yourdomain.tld`: Your root domain.

---

## **How It Works**
- **Manual Auth Hook**: The plugin creates the required `_acme-challenge` TXT records in Vultr's DNS for domain validation.
- **Manual Cleanup Hook**: After validation, the plugin cleans up the TXT records to maintain a clean DNS configuration.

---

## **Contributing**
We welcome contributions! Feel free to open an issue or submit a pull request if you find any bugs or want to add features.

---

## **License**
This project is licensed under the Do Whatever You Want License.
