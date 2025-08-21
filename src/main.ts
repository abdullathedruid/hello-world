// htmx is already loaded via script tag, so we can use it directly

// Your TypeScript code goes here
console.log('Hello from TypeScript!');

// Example: Add custom htmx behavior
document.addEventListener('DOMContentLoaded', () => {
  if (window.htmx) {
    window.htmx.on('htmx:afterRequest', (event: any) => {
      console.log('HTMX request completed:', event.detail);
    });
  }
});