import React from 'react';
import './Board.css';
import Thread from "./Thread";
import {Link} from "react-router-dom";
import OpeningPost from "./OpeningPost";

class Board extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      time: this.props.time,
      threads: []
    }
  }

  getAllThreads() {
    let worldTime = ""
    if (this.state.time !== null && typeof this.state.time !== 'string') {
      worldTime = this.state.time.unix()
    }

    fetch(this.props.apiUrl + '/thread/all?time=' + worldTime)
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

  componentDidUpdate(prevProps, prevState, snapshot) {
    if (this.props.name !== prevProps.name) {
      this.getAllThreads()
    }
    return true
  }

  render() {
    return (
        <div>
          {this.displayNav()}
          <h2>/{this.props.name}/</h2>

          {this.allThreads()}
        </div>
    );
  }

  displayNav() {
    if (this.props.catalog) {
      return <span className={"timer"}>
            <Link to={"/" + this.props.name + "/"}>Return</Link>
          </span>;
    }

    return <span className={"timer"}>
            <Link to={"/" + this.props.name + "/catalog"}>Catalog</Link>
          </span>;
  }

  allThreads() {
    if (this.state.threads == null) {
      return (<div> . . . </div>)
    }
    if (this.state.threads.length === 0) {
      return (<div> . </div>)
    }

    if (this.props.catalog) {
      return <div className={"grid-container"}>{this.state.threads.map(thread => {
        return (
              <OpeningPost key={thread.post.no} board={this.props.name} thread={thread} imageContext={this.props.imageContext}/>)
      })}
      </div>
    }

    let threads = this.state.threads.slice(0, 10);
    return threads.map((thread) => {

      return (<Thread key={thread.post.no} board={this.props.name}
                      apiUrl={this.props.apiUrl}
                      imageContext={this.props.imageContext}
                      thread={thread}
                      limit='5'/>);
    })
  }

}

export default Board