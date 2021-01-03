import React from 'react';
import Datetime from 'react-datetime'

function TimeForm(props) {

    let time = props.time
    let setTime = props.timeSetter


    return (
        <div>
            <Datetime
                value={time}
                onChange={setTime}
            />
        </div>
    )
}

export default TimeForm