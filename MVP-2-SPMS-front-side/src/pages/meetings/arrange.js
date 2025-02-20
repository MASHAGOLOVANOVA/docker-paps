import React, { useState, useEffect, useRef } from 'react';
import Cookies from 'universal-cookie';
import { Row, Button, Col,Form, Modal } from 'react-bootstrap';

import SpinnerCenter from 'pages/shared/Spinner';
import { GetUtcDate } from 'pages/shared/FormatDates';

const cookies = new Cookies();

const getReqOptions = {
    method: "GET",
    mode: "cors",
    cache: "default",
    credentials: 'include',
    headers: {
        "Session-Id": cookies.get('session_token')
    }
};

let postMeetingReqOptions = {
    method: "POST",
    mode: "cors",
    cache: "default",
    credentials: 'include',
    headers: {
        "Session-Id": cookies.get('session_token'),
        "Content-Type": "application/json",
    },
};

function ArrangeMeeting() {
    const [students, setStudents] = useState(null);
    const [meetings, setMeetings] = useState(null);
    const [formData, setFormData] = useState({});
    const [showAddMeetingResult, setShowAddMeetingResult] = useState(false);
    const [addMeetingResult, setAddMeetingResult] = useState(null);
    const [integr, setIntegr] = useState(null);

    const selectRef = useRef(null);
    const timeRef = useRef(null);

    useEffect(() => {

        // setting meeting defaults
        setFormData({
            "name": "",
            "description": "",
            "meeting_time": "",
            "student_participant_id": -1,
            "is_online": false
        });
        fetch(`${process.env.REACT_APP_SERVER_ADDR}/api/v1/account/integrations`, getReqOptions)
            .then(response => response.json())
            .then(json => {
                setIntegr(json);
                if (json.planner) {
                    fetch(`${process.env.REACT_APP_SERVER_ADDR}/api/v1/students`, getReqOptions)
                        .then(response => response.json())
                        .then(json => { setStudents(json["students"]); resetSelect(); })
                        .catch(error => console.error(error));
                }
            });
    }, []);

    function resetSelect() {
        // set select as invalid, so student has to be chosen
        const element = selectRef.current;
        element.setCustomValidity("Выберите студента!");
    }

    function resetTime() {
        const element = timeRef.current;
        element.value = "";
    }
    const handleChange = (event) => {
        const name = event.target.name;
        const value = event.target.value;
        setFormData(values => ({ ...values, [name]: value }));
    }

    const handleCehckboxChange = (event) => {
        const name = event.target.name;
        const checked = event.target.checked;
        setFormData(values => ({ ...values, [name]: checked }));
    }

    const handleSelectChange = (event) => {
        // you cannot select default option anymore, so we can just make it valid
        const vals = event.target.value.split(";");
        const proj = parseInt(vals[1]);
        const stud = parseInt(vals[0]);
        setFormData(values => ({ ...values, ["student_participant_id"]: stud }))
        setFormData(values => ({ ...values, ["project_id"]: proj }))
        event.target.setCustomValidity("");
    }

    async function handleSubmit(event) {
        event.preventDefault();
        resetTime();
        prepareReqBody();
        OpenRequestResultModal()
        try {
            let to = GetUtcDate(formData["meeting_time"]);
            to.setTime(to.getTime() + 60 * 60000 - to.getTimezoneOffset() * 60 * 1000)
            let from = GetUtcDate(formData["meeting_time"]);
            from.setTime(from.getTime() - 60 * 60000 - from.getTimezoneOffset() * 60 * 1000)
            const meetingsInTimeRes = await fetch(`${process.env.REACT_APP_SERVER_ADDR}/api/v1/meetings?from=${from.toISOString()}&to=${to.toISOString()}`, getReqOptions)
            const meetingsjs = await meetingsInTimeRes.json()
            setMeetings(meetingsjs["meetings"])

            if (meetingsjs["meetings"].length === 0) {
                const response = await fetch(`${process.env.REACT_APP_SERVER_ADDR}/api/v1/meetings/add`, postMeetingReqOptions)
                const status = response.status;
                console.log("Responce status:", status);
                event.target.reset();
                resetSelect();
                setAddMeetingResult(status);
                return
            }
            setAddMeetingResult(-1)
        } catch (error) {
            console.error("Error:", error);
        }
    }

    function prepareReqBody() {
        formData["meeting_time"] += ":00.000Z";
        formData["student_participant_id"] = parseInt(formData["student_participant_id"])
        postMeetingReqOptions["body"] = JSON.stringify(formData)
    }

    function RenderStudents() {
        return students.map((student) =>
            <option value={`${student.id};${student.project_id}`}>
                {`${student.surname} ${student.name[0]}. ${student.middlename[0]}., ${student.cource} курс ${student.education_programme}, ${student.project_theme}`}
            </option>)
    }

    function OpenRequestResultModal() {
        setShowAddMeetingResult(true);
    }
    function CloseRequestResultModal() {
        setShowAddMeetingResult(false);
        setAddMeetingResult(null);
    }

    function RenderRequestResultModal() {
        let header = "Встреча назначена!";
        let body = "Вы можете просмотреть её в своем расписании здесь или в подключенном календаре.";
        if (addMeetingResult !== 200) {
            if (addMeetingResult === -1) {
                header = "Пересечение встреч!";
                body = `На данное время уже назначена встреча! Она начинается в: ${meetings[0].time}`;
            }
            else {
                header = "Произошла ошибка при назанчении встречи!";
                body = `Код ошибки: ${addMeetingResult}. Обратитесь в службу поддержки, если прблема не устранится.`;
            }
        }
        return <>
            <Modal.Header>
                <Modal.Title>{header}</Modal.Title>
            </Modal.Header>
            <Modal.Body>{body}</Modal.Body>
            <Modal.Footer>
                <Button className='style-button' onClick={CloseRequestResultModal}>
                    ОК
                </Button>
            </Modal.Footer>
        </>
    }

    return (
        <>
            <Row className='justify-content-center'>
                <Col xs={11} md={10} lg={8}>
                    <h1 className='mb-4'>Назначить встречу</h1>
                    <hr />
                    <div >
                        {integr ? integr.planner ?
                            <Row className='justify-content-center'>
                                <Col xs={12} sm={8}>
                                    <Form onSubmit={handleSubmit}>
                                        <Form.Group className="mb-3" controlId="meetName">
                                            <Form.Label>Название *</Form.Label>
                                            <Form.Control name="name" onChange={handleChange} required placeholder="Введите название" />
                                            <Form.Text className="text-muted">
                                                Будьте кратки в названии
                                            </Form.Text>
                                        </Form.Group>

                                        <Form.Group className="mb-3" controlId="meetDesc">
                                            <Form.Label>Описание</Form.Label>
                                            <Form.Control name="description" onChange={handleChange} placeholder="Введите описание" />
                                        </Form.Group>

                                        <Form.Group className="mb-3" controlId="meetTime">
                                            <Form.Label>Дата и время *</Form.Label>
                                            <Form.Control ref={timeRef} name="meeting_time" onChange={handleChange} required type="datetime-local"
                                                placeholder="Введите время встречи"
                                                id="meeting-time"
                                                min={new Date(Date.now()).toISOString().split(":")[0] + ":" + new Date(Date.now()).toISOString().split(":")[1]}
                                            />
                                        </Form.Group>

                                        <Form.Group className="mb-3" controlId="student">
                                            <Form.Label>С кем встреча *</Form.Label>
                                            <Form.Select ref={selectRef} name="-" onChange={handleSelectChange} required aria-label="select student" >
                                                <option value={-1} selected hidden>Выберите студента...</option>
                                                {students ? RenderStudents() :
                                                    <option disabled>{SpinnerCenter()}</option>
                                                }
                                            </Form.Select>
                                        </Form.Group>
                                        <Form.Group className="mb-3 " controlId="isOnline">
                                            <label class="style-checkmark-label">
                                                <input name="is_online" onChange={handleCehckboxChange} type="checkbox" class="style-default-checkmark" />
                                                <span class="style-checkmark"></span>
                                                Онлайн
                                            </label>
                                        </Form.Group>
                                        <Button type="submit" className="style-button">
                                            Назначить
                                        </Button>
                                    </Form>
                                </Col>
                            </Row>
                            :
                            <>
                                <h3>Вы еще не подключили планировщик, это можно сделать <a href='/profile/integrations'>здесь</a></h3>
                            </> :
                            SpinnerCenter()}

                    </div>
                </Col>
            </Row>
            <Modal backdrop="static" show={showAddMeetingResult} onHide={CloseRequestResultModal}>
                {addMeetingResult ? RenderRequestResultModal() :
                    <Modal.Header className="justify-content-md-center">
                        {SpinnerCenter()}
                    </Modal.Header>}
            </Modal>
        </>
    );
};

export default ArrangeMeeting;