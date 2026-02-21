import React, { useState, useEffect } from "react";
import Button from "../shared/Button";

const STREAM_TOKEN = "streamToken";

interface Token {
  token: string;
}

const StreamTokenView = () => {
  const [token, setToken] = useState<Token>();
  const [showToken, setShowToken] = useState<boolean>(false);

  const fetchToken = () => {
    fetch("/api/profiles/get", {
      method: "GET",
    })
      .then((result) => {
        if (result.status == 200) {
          console.log(result);
          return result.json();
        } else if (result.status == 403) {
          setToken({ token: "Forbidden" });
          return;
        }
        setToken({ token: "Unhandled error retrieving stream token" });
        return;
      })
      .then((result) => {
        if (result) {
          localStorage.setItem(STREAM_TOKEN, result["token"]);
          setToken(() => result);
        }
      })
      .catch((error) => {
        console.log(error);
        setToken({
          token: "Error processing response",
        });
      });
  };

  const resetToken = () => {
    fetch("/api/profiles/reset", {
      method: "POST",
    }).then((result) => {
      if (result.status == 200) {
        console.log("Reset token");
        fetchToken();
      } else {
        console.log("Failed to reset token");
      }
    });
  };

  const toggleShowToken = () => {
    setShowToken(!showToken);
  };

  const copyTokenToClipboard = () => {
    if (token) {
      navigator.clipboard.writeText(token.token);
    }
  };

  useEffect(() => {
    var token = localStorage.getItem(STREAM_TOKEN);
    if (token != null) {
      setToken({
        token: token,
      });
    } else {
      fetchToken();
    }
  }, []);

  return (
    <div className="flex flex-col">
      <div>Stream token</div>
      <div className="flex flex-row max-w-2xl pt-1 gap-2 justify-center">
        <Button title={"Reset"} center onClick={resetToken} />
        <Button title={"Refresh"} center onClick={fetchToken} />
        <Button title={"Copy"} center onClick={copyTokenToClipboard} />
        <Button
          title={showToken === true ? "Hide" : "Show"}
          center
          onClick={toggleShowToken}
        />
      </div>
      <div>
        <div className="flex flex-row pt-5 gap-2 justify-center">
          {showToken === true ? token?.token : "TOKEN HIDDEN"}
        </div>
      </div>
    </div>
  );
};

export default StreamTokenView;
