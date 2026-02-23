// Select DOM elements
const urlInput = document.getElementById('urlInput');
const shortenBtn = document.getElementById('shortenBtn');
const resultContainer = document.getElementById('resultContainer');
const shortUrlElement = document.getElementById('shortUrl');
const copyBtn = document.getElementById('copyBtn');
const copyMessage = document.getElementById('copyMessage');

// Event Listener for the Shorten Button
shortenBtn.addEventListener('click', function() {
    const longUrl = urlInput.value.trim();

    // Basic validation
    if (longUrl === '') {
        alert('Please paste a URL first!');
        return;
    }

    // Generate a random 6-character code
    const randomCode = generateRandomCode(6);
    
    // Create the fake short URL
    const shortUrl = `NIET.ly/${randomCode}`;

    // Update the DOM
    shortUrlElement.textContent = shortUrl;
    shortUrlElement.href = longUrl; // In a real app, this would point to the redirect logic
    
    // Show the result section
    resultContainer.classList.remove('hidden');
    
    // Reset input field
    urlInput.value = '';
});

// Event Listener for the Copy Button
copyBtn.addEventListener('click', function() {
    const urlText = shortUrlElement.textContent;

    // Use the Clipboard API to copy text
    navigator.clipboard.writeText(urlText).then(() => {
        // Show success message
        showCopyMessage();
    }).catch(err => {
        console.error('Failed to copy: ', err);
    });
});

// Helper Function: Generate Random String
function generateRandomCode(length) {
    const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    let result = '';
    for (let i = 0; i < length; i++) {
        result += characters.charAt(Math.floor(Math.random() * characters.length));
    }
    return result;
}

// Helper Function: Show "Copied" message temporarily
function showCopyMessage() {
    copyMessage.classList.add('show');
    
    // Hide the message after 2 seconds
    setTimeout(() => {
        copyMessage.classList.remove('show');
    }, 2000);
}