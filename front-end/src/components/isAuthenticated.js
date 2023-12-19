function isAuthenticated() {
  return fetch('http://localhost:8080/api/validateSession', {
    method: 'GET',
    credentials: 'include', // Include cookies in the request
  })
    .then((response) => {
      if (response.ok) {
        return true; // User is authenticated
      } else {
        return false;
      }
    })
    .catch((error) => {
      console.error('Authentication error:', error);
      return false;
    });
}

export default isAuthenticated;
