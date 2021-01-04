import React from 'react';
import {Link} from "react-router-dom";

import './Board.css';
import './Hover.css';

import objection from './objection.gif'
import TimeForm from "./TimeForm";
import {Hover} from "./Hover";


class Thread extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            status: "SUCCESS",
            board: this.props.board,
            apiUrl: this.props.apiUrl,
            imageContext: this.props.imageContext,
            no: this.props.no,
            thread: this.props.thread,
            limit: this.props.limit
        }
    }

    componentDidMount() {
        if (this.state.thread === undefined || this.state.thread === null || this.state.status === "FAILURE") {
            this.getThread()
        }
    }

    getThread() {
        fetch(this.state.apiUrl + '/thread?no=' + this.state.no)
            .then((response) => {
                return response.json()
            })
            .then((json) => {
                console.log(json);
                this.setState({
                    status: json.status,
                    thread: json.thread
                })
            }).catch(console.log);
    }

    displayImage(post) {
        if (this.state.imageContext.startsWith("http")) {
            return this.state.imageContext + "/" + post.image
        }
        return process.env.PUBLIC_URL + this.state.imageContext + post.image;
    }

    displayThread(thread) {
        return (
            <div key={thread.post.no}>
                <hr/>
                <div className="thread">
                    <span className="image"><img alt={thread.post.filename}
                                                 src={this.displayImage(thread.post)}/></span><span
                    className="threadHeader">{thread.subject} <span
                    className="postName">{thread.post.name}</span> {thread.post.timestamp} No. <Link
                    to={"/" + this.state.board + "/res/" + thread.post.no}>{thread.post.no}</Link> <span
                    className="quotedBy">{this.quotedBy(thread.post, thread)}</span></span>

                    <div><span className="content">{this.displayComment(thread.post, thread)}</span></div>
                </div>
                <div className="replies">
                    {this.displayReplies(thread, 5)}
                </div>
            </div>
        );
    }

    displayReplies(thread) {
        if (this.state.limit !== null) {
            thread.replies = thread.replies.slice(-this.state.limit)
        }
        return thread.replies.map((post) => {
            return (
                this.displayPost(post, thread)
            )
        })
    }

    displayPost(post, thread, hover = false) {
        return <div key={post.no} className="post">
            {this.optionalImage(post)}
            <span className="postHeader"><span
                className="postName">{post.name}</span> {post.timestamp} No. {post.no} <span
                className="quotedBy">{this.quotedBy(post, thread, hover)}</span></span>
            <div><span className="content">{this.displayComment(post, thread)}</span></div>
        </div>;
    }

    displayComment(post, thread) {
        if (post.comment_segments == null) {
            return post.comment;
        }

        return post.comment_segments.map((segment, i) => {
            return this.displaySegment(segment, i, post, thread)
        })
    }

    displaySegment(segment, i, post, thread) {
        switch (segment.format[0]) {
            case "objection":
                return (
                    <div className={segment.format}><img src={objection} alt="Objection!"/></div>
                );
            case "roll":
                return (
                    <div className="roll">{Math.random() % 6}</div>
                );
            case "noQuote":
                return (
                    <div className={this.formatAsClasses(segment.format)} key={i}>
                        <Hover key={i} onHover={this.displayPostHover(this.findPost(thread, parseInt(segment.segment.replace(">>", ""))), thread)}>
                            <span className="noQuote" key={i}>{segment.segment}</span>
                        </Hover>
                    </div>
                );
            default:
                return (
                    <div className={this.formatAsClasses(segment.format)} key={i}>{segment.segment}<br/></div>
                );
        }
    }

    formatAsClasses(segment) {
        if (segment === null || segment.format === null) {
            return "";
        }
        return segment.map((format) => {
            return format + " "
        })
    }

    quotedBy(post, thread, hover = false) {
      if (post.quoted_by == null || hover) {
            return <span/>
        }
        return post.quoted_by.map((postNo, i) => {
            return (
                <Hover key={i} onHover={this.displayPostHover(this.findPost(thread, postNo), thread)}>
                    <span className="noQuote" key={i}>>>{postNo} </span>
                </Hover>

            )
        })
    }

    displayPostHover(post, thread) {
        if (post === null) {
            return <span className="hoveredPost"/>
        }

        return this.displayPost(post, thread, true)
    }

    optionalImage(post) {
        if (post.image != null && post.image !== "") {
            return <span className="image"><img src={this.displayImage(post)}
                                                alt={post.filename}/></span>
        } else {
            return <span className="noImage"/>
        }
    }

    render() {
        if (this.state.thread === undefined || this.state.thread === null || this.state.status === "FAILURE") {
            return (<div>. . .</div>)
        }

        return (
            <div>
                {this.displayTimer()}
                {this.displayThread(this.state.thread)}
            </div>
        )
    }

    findPost(thread, postNo) {
        if (thread.post.no === postNo) {
            return thread.post
        }

        for (let i = 0; i < thread.replies.length; i++) {
            if (thread.replies[i].no === postNo) {
                return thread.replies[i]
            }
        }
        return null;
    }

    displayTimer() {
        let enabled = this.state.limit === undefined || this.state.limit === 0 || this.state.limit === null
        if (enabled) {
            return <TimeForm time={new Date()} timeSetter={(v) => {}}/>
        }
    }
}

export default Thread