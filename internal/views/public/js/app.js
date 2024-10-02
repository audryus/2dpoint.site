const cookies = document.cookie.split(';');
console.log(cookies)
const desiredCookie = cookies.find(cookie => cookie.includes('2dpoint_token='));
console.log(desiredCookie)
if (desiredCookie) {
    const [_, value] = desiredCookie.split('=');
    document.body.setAttribute('hx-headers', `{"X-Csrf-Token": "${value}"}`);
}

