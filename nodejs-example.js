const axios = require("axios");

const data = {
  key: "testJSON",
  value: { hello: "world" }, // JSON value
};

axios
  .post("http://localhost:8080/set", data)
  .then((response) => {
    console.log(response.data); // Log the response from the server

    axios
      .get(`http://localhost:8080/get?key=${data.key}`)
      .then((response) => {
        console.log(response.data); // Log the response from the server
      })
      .catch((error) => {
        console.error("Error:", error);
      });
  })
  .catch((error) => {
    console.error("Error:", error.code);
    console.error(error.message);
  });
