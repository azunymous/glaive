import React, {useEffect, useState} from 'react';
import {
  BrowserRouter as Router,
  Link,
  Route,
  Switch,
  useParams
} from "react-router-dom";

import Board from './Board'
import Thread from "./Thread";
import './Board.css';
import moment from "moment";
import Datetime from "react-datetime";

function App() {
  let [boards, setBoards] = useState([]);

  function getTimestampFromStorageOrDefault() {
      let storedTimestamp = localStorage.getItem('timestamp')
      if (storedTimestamp === null) {
          return moment(1340402891 * 1000)
      }
      return moment(parseInt(storedTimestamp))
  }

  let [worldTime, setWorldTime] = useState(getTimestampFromStorageOrDefault);
  function setWorldTimeAndLocalStorage(time) {
      localStorage.setItem('timestamp', (time.unix() * 1000).toString())
      setWorldTime(time)
  }

  useEffect(getBoards, []);

  return (
      <div className="App">
        <header className="App-header">
          <h1>igiari.net</h1>
        </header>
        <div className="outer">
          <Router>
                    <span>
                        <ul id="menu">
                            <li><Link to="/">Home</Link></li>
                        <ShowBoardLinks/>
                        </ul>
                    </span>
            <Switch>
              <Route exact path="/">
                . . .
              </Route>
              <Route exact path="/:boardID/" children={<ShowBoard/>}/>
              <Route path="/:boardID/res/:threadNo" children={<ShowThread/>}/>
            </Switch>
          </Router>
        </div>
      </div>
  );

  function getBoards() {
    fetch(process.env.REACT_APP_API_URL + '/boards')
    .then((response) => {
      return response.json()
    })
    .then((json) => {
      console.log(json);
      setBoards(json);
    }).catch(console.log);
  }

  function getBoardDetails(boardID) {
    let board = boards["/" + boardID + "/"];
    let apiURL = board["host"];
    let imageContext = board["images"];
    return {apiURL, imageContext};
  }

  function ShowBoard() {
    let {boardID} = useParams();
    if (boards.length === 0) {
      return (
          <div>
            Loading...
          </div>
      )
    }
    let {apiURL, imageContext} = getBoardDetails(boardID);
    return (
        <div>
            <WorldClock time={worldTime} timeSetter={setWorldTime}/>
            <Board name={boardID} apiUrl={apiURL} imageContext={imageContext} time={worldTime}/>
        </div>

    );
  }

  function ShowThread() {
    let {boardID, threadNo} = useParams();
    if (boards.length === 0) {
      return (
          <div>
            Loading...
          </div>
      )
    }
    let {apiURL, imageContext} = getBoardDetails(boardID);
    return (
        <div>
            <WorldClock time={worldTime} timeSetter={setWorldTime}/>
            <Thread board={boardID} no={threadNo} apiUrl={apiURL} imageContext={imageContext} time={worldTime} setTime={setWorldTimeAndLocalStorage}/>
        </div>
    );
  }

  function ShowBoardLinks() {
    if (boards === undefined || boards === null) {
      return <li>No Boards Found</li>
    }

    return Object.keys(boards).map((id) => {
      return (
          <li key={id}><Link to={id}>{id}</Link></li>
      )
    })
  }


    function WorldClock(props) {
        let time = props.time
        return (
            <div>
                <Datetime
                    value={time}
                    onChange={setWorldTimeAndLocalStorage}
                />
            </div>
        )
    }
}




export default App;
