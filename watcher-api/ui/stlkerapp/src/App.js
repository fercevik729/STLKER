import './App.css';
import React, {useState, useEffect} from 'react';
import { WatcherClient } from './proto/watcher_grpc_web_pb';
import { TickerRequst} from './proto/watcher_pb';

// Create a client to connect to API

var client = new WatcherClient("https://localhost:8080")
function App() {
  // Create a constant named status and a func called setStatus
  const sendTicker = () => {
    // sendReq is a func that will send a ticker request to the backend
    var tickerRequest = new TickerRequst();
    // use the client to send the ticker request
    client.echo(tickerRequest, null, function(err, response) {
      var msg = response.toObject();
      console.log(msg.getTicker());
      // call set status to change the value of status
    })
  }

  return (
    <div className="App">
      <p>Look at the log</p>
    </div>
  );
}

export default App;
