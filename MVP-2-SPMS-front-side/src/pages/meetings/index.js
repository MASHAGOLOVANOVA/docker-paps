import React, { useState, useEffect } from 'react';
import Cookies from 'universal-cookie';
import {  Row,  Badge, Col, Alert} from 'react-bootstrap';
import LinkContainer from 'react-router-bootstrap/LinkContainer';
import { Link } from 'react-router-dom';
import SpinnerCenter from 'pages/shared/Spinner';
import { FormatDate, FormatDateTime } from 'pages/shared/FormatDates';

const cookies = new Cookies();

const reqOptions = {
    method: "GET",
    mode: "cors",
    cache: "default",
    credentials: 'include',
    headers: {
        "Session-Id": cookies.get('session_token')
    }
};

function Meetings() {
    const [meetings, setMeetings] = useState(null);
    const [meetingsStatuses, setMeetingsStatuses] = useState(null);
    const [meetingFromTime, setMeetingFromTime] = useState(new Date())
    const [integr, setIntegr] = useState(null);


    useEffect(() => {
        fetch(`${process.env.REACT_APP_SERVER_ADDR}/api/v1/account/integrations`, reqOptions)
            .then(response => response.json())
            .then(json => {
                setIntegr(json);
                if (json.planner) {
                    const currentTime = new Date(Date.now());
                    fetch(`${process.env.REACT_APP_SERVER_ADDR}/api/v1/meetings?` + new URLSearchParams({
                        from: currentTime.toISOString(),
                    }), reqOptions)
                        .then(response => response.json())
                        .then(json => setMeetings(json["meetings"]))
                        .catch(error => console.error(error));
                    fetch(`${process.env.REACT_APP_SERVER_ADDR}/api/v1/meetings/statuslist`, reqOptions)
                        .then(response => response.json())
                        .then(json => setMeetingsStatuses(json["statuses"]))
                        .catch(error => console.error(error));
                    setMeetingFromTime(currentTime)
                    return
                }
            })
            .catch(error => console.error(error));


    }, []);

    function ParseMeetings() {
        if (meetings.length === 0) {
            return (
                <div>
                    Встреч нет!
                </div>
            )
        }
        return meetings.map((meeting) => <>
            <Row as="li" className="d-flex justify-content-between align-items-start ">
                <Col xs={12} sm={12} md={3} lg={2} xxl="auto" className='fst-italic text-secondary'>
                    {FormatDateTime(meeting.time)}
                </Col>
                <Col className='mb-3 text-break'>
                    <h5>
                        {meeting.name} <Badge pill className='style-bg'>{meeting.is_online ? "online" : "offline"}</Badge> 
                    </h5>
                    <div className=''>
                        Студент: <LinkContainer to={"/projects/" + meeting.student.project_id} className='fst-italic' title={meeting.student.project_theme}>
                            <Link className="link-body-emphasis link-offset-2 link-underline-opacity-0 link-underline-opacity-50-hover">
                                {meeting.student.name}, {meeting.student.cource} курс
                            </Link>
                        </LinkContainer>
                    </div>
                    <div className='text-secondary-emphasis'>
                        {meeting.description}
                    </div>
                    {meeting.has_planner_meeting ? <></> :
                        <div>
                            <Alert variant="danger" className='my-3'>
                                У данной встречи нет события в календаре!
                            </Alert>
                        </div>
                    }
                </Col>
                <hr />
            </Row>
        </>
        )
    }

    return (
        <>
            <Row className='justify-content-center'>
                <Col xs={11} md={10} lg={8}>
                    <h1 className='mb-4'>Встречи с {FormatDate(meetingFromTime.toISOString())}</h1>
                    <hr />
                    <div >
                        {integr ?
                            integr.planner ?
                                meetings && meetingsStatuses? ParseMeetings() : SpinnerCenter() :
                                <>
                                    <h3>Вы еще не подключили планировщик, это можно сделать <a href='/profile/integrations'>здесь</a></h3>
                                </> :
                            SpinnerCenter()
                        }

                    </div>
                </Col>
            </Row>
        </>
    );
};

export default Meetings;