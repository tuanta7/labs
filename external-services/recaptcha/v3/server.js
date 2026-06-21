require("dotenv").config();

const express = require("express");
const axios = require("axios");
const cors = require("cors");
const app = express();
const port = 3000;

app.use(express.json());
app.use(cors());

const SECRET_KEY = process.env.RECAPTCHA_SECRET_KEY;

app.post("/verify", async (req, res) => {
  try {
    const { token } = req.body;

    if (!token) {
      return res.status(400).json({ error: "Token is required" });
    }

    const params = {
      secret: SECRET_KEY,
      response: token,
    };

    console.log(params);

    const response = await axios.post(
      "https://www.google.com/recaptcha/api/siteverify",
      null,
      {
        params,
      }
    );

    const { success, score, action } = response.data;
    console.log(response.data);

    if (success) {
      res.json({
        success: true,
        score,
        action,
        message: "reCAPTCHA verification successful",
      });
    } else {
      res.status(400).json({
        success: false,
        message: "reCAPTCHA verification failed",
      });
    }
  } catch (error) {
    console.error("Error verifying reCAPTCHA:", error);
    res.status(500).json({
      error: "Server error",
      message: error.message,
    });
  }
});

app.listen(port, () => {
  console.log(`Example app listening on port ${port}`);
});
