import React from 'react';
import './Board.css';
import Thread from "./Thread";
import TimeForm from "./TimeForm";
import moment from "moment";

class Board extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      name: this.props.name,
      apiUrl: this.props.apiUrl,
      imageContext: this.props.imageContext,
      time: moment(1340402891 * 1000),
      threads: []
    }
    this.setTime = this.setTime.bind(this)
  }

  setTime(time) {
    console.log("Setting time to " + time)
    this.setState({
      time: time
    })
    this.getAllThreads()
  }


  getAllThreads() {
    let time = ""
    if (this.state.time !== null && typeof this.state.time !== 'string') {
      time = this.state.time.unix()
    }
    fetch(this.state.apiUrl + '/thread/all?time=' + time)
    .then((response) => {
      return response.json()
    })
    .then((json) => {
      console.log(json);
      this.setState({
        threads: json
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
    let threads = this.state.threads.slice(0, 10);
    return threads.map((thread) => {
      return (<Thread key={thread.post.no} board={this.state.name}
                      apiUrl={this.state.apiUrl}
                      imageContext={this.state.imageContext} thread={thread}
                      limit='5'/>);
    })
  }

}

export default Board