const cookies = document.cookie.split(';');
const desiredCookie = cookies.find(cookie => cookie.includes('2dpoint_token='));

if (desiredCookie) {
    const [_, value] = desiredCookie.split('=');
    document.body.setAttribute('hx-headers', `{"X-Csrf-Token": "${value}"}`);
}

