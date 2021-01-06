import React from 'react';
import './Board.css';
import Thread from "./Thread";
import TimeForm from "./TimeForm";
import moment from "moment";

class Board extends React.Component {
  constructor(props) {
    super(props);

    function getTimestampFromStorageOrDefault() {
      let storedTimestamp = localStorage.getItem('timestamp')
      if (storedTimestamp === null) {
        return moment(1340402891 * 1000)
      }
      return moment(parseInt(storedTimestamp))
    }

    this.state = {
      name: this.props.name,
      apiUrl: this.props.apiUrl,
      imageContext: this.props.imageContext,
      time: getTimestampFromStorageOrDefault(),
      threads: []
    }
    this.setTime = this.setTime.bind(this)
  }

  setTime(time) {
    console.log("Setting time to " + time)
    this.setState({
      time: time
    }, () => this.getAllThreads())
    localStorage.setItem('timestamp', (time.unix() * 1000).toString())
  }


  getAllThreads() {
    let worldTime = ""
    if (this.state.time !== null && typeof this.state.time !== 'string') {
      worldTime = this.state.time.unix()
    }

    fetch(this.state.apiUrl + '/thread/all?time=' + worldTime)
    .then((response) => {
      return response.json()
    })
    .then((json) => {
      if (json.status !== null && json.status === "FAILURE") {
        console.log(json)
        return
      }

      this.setState({
        threads: json,
      })
    }).catch(console.log);
  }

  componentDidMount() {
    this.getAllThreads();
  }

  render() {
    return (
        <div>
          <h2>/{this.state.name}/</h2>
          <TimeForm apiUrl={this.state.apiUrl} time={this.state.time} timeSetter={this.setTime}/>
          {this.allThreads()}
        </div>
    );
  }

  allThreads() {
    if (this.state.threads == null) {
      return (<div> . . . </div>)
    }
    if (this.state.threads.length === 0) {
      return (<div> . </div>)
    }

    console.log(this.state.threads)
    let threads = this.state.threads.slice(0, 10);
    return threads.map((thread) => {

      return (<Thread key={thread.post.no} board={this.state.name}
                      apiUrl={this.state.apiUrl}
                      imageContext={this.state.imageContext}
                      thread={thread}
                      limit='5'/>);
    })
  }

}

export default Board