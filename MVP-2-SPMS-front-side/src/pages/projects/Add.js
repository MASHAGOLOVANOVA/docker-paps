import { Row, Button, Col, Form, Tab, Nav, Modal} from 'react-bootstrap';
import Cookies from 'universal-cookie';
import React, { useState, useEffect, useRef } from 'react';
import SpinnerCenter from 'pages/shared/Spinner';

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

let postReqOptions = {
    method: "POST",
    mode: "cors",
    cache: "default",
    credentials: 'include',
    headers: {
        "Session-Id": cookies.get('session_token'),
        "Content-Type": "application/json",
    },
};

function AddProject() {

    const steps = ["select-student-method", "add-student", "select-student", "add-project"];

    const [students, setStudents] = useState(null);
    const [edprogs, setEdprogs] = useState(null);
    const [currentStepIndex, setCurrentStep] = useState(0);
    const [navbar, setNavbar] = useState(<></>);
    const [formData, setFormData] = useState({});
    const [studentFormData, setStudentFormData] = useState({});
    const [addingStudent, setAddingStudent] = useState(false);
    const [isStudentSelected, setIsStudentSelected] = useState(false);
    const [isStudentFilled, setIsStudentFilled] = useState(false);
    const [showAddProjectResult, setShowAddProjectesult] = useState(false);
    const [addProjectResult, setAddProjectResult] = useState(null);

    const methodStudentStep = 0;
    const addStudentStep = 1;
    const selectStudentStep = 2;
    const addProjectStep = 3;

    const [stepPath, setPath] = useState([methodStudentStep, selectStudentStep, addProjectStep]);
    const [integr, setIntegr] = useState(null);

    const selectRef = useRef(null);
    const addStudentRef = useRef(null);

    useEffect(() => {
        // setting project defaults
        setFormData({
            "theme": "",
            "student_id": 1,
            "year": 2024,
            "repository_owner_login": "",
            "repository_name": ""
        });
        // setting student defaults
        setStudentFormData({
            "name": "",
            "surname": "",
            "middlename": "",
            "cource": 1,
            "education_programme_id": 1
        });

        fetch(`${process.env.REACT_APP_SERVER_ADDR}/api/v1/account/integrations`, getReqOptions)
            .then(response => response.json())
            .then(json => {
                setIntegr(json);
                if (json.cloud_drive) {
                    fetch(`${process.env.REACT_APP_SERVER_ADDR}/api/v1/students`, getReqOptions)
                        .then(response => response.json())
                        .then(json => setStudents(json["students"]))
                        .catch(error => console.error(error));
                    fetch(`${process.env.REACT_APP_SERVER_ADDR}/api/v1/universities/1/edprogrammes`, getReqOptions)
                        .then(response => response.json())
                        .then(json => setEdprogs(json["programmes"]))
                        .catch(error => console.error(error));
                    resetSelect();
                    renderNavBar();
                }
            });

    }, []);


    useEffect(() => {
        renderNavBar();
    }, [currentStepIndex, stepPath, isStudentFilled, isStudentSelected]);


    function resetSelect() {
        // set select as invalid, so student has to be chosen
        const element = selectRef.current;
        setIsStudentSelected(false)
        if (element) {
            element.setCustomValidity("Выберите студента!");
        }
    }

    const handleChange = (event) => {
        const name = event.target.name;
        const value = event.target.value;
        setFormData(values => ({ ...values, [name]: value }));
    }

    const handleStudentChange = (event) => {
        const name = event.target.name;
        const value = event.target.value;
        setStudentFormData(values => ({ ...values, [name]: value }));
        checkStudentFormValidity();
    }

    function checkStudentFormValidity() {
        const element = addStudentRef.current;
        setIsStudentFilled(element.checkValidity());
    }

    const handleSelectChange = (event) => {
        // you cannot select default option anymore, so we can just make it valid
        event.target.setCustomValidity("");
        setIsStudentSelected(true);
        handleChange(event);
    }

    async function handleSubmit(event) {
        event.preventDefault();
        event.target.reset();
        resetSelect();
        OpenRequestResultModal()
        try {
            if (addingStudent) {
                prepareStudentReqBody()
                const response = await fetch(`${process.env.REACT_APP_SERVER_ADDR}/api/v1/students/add`, postReqOptions)
                if (response.status !== 200) {
                    const status = response.status;
                    setAddProjectResult(status)
                    return
                }
                const data = await response.json();
                formData["student_id"] = data["student_id"];
            }
            prepareProjectReqBody();
            const response = await fetch(`${process.env.REACT_APP_SERVER_ADDR}/api/v1/projects/add`, postReqOptions)
            const status = response.status;
            console.log("Responce status:", status);
            setAddProjectResult(status)
        } catch (error) {
            console.error("Error:", error);
        }
    }

    function prepareProjectReqBody() {
        if (typeof formData["student_id"] === 'string') {
            formData["student_id"] = parseInt(formData["student_id"])
        }
        formData["year"] = parseInt(formData["year"])
        postReqOptions["body"] = JSON.stringify(formData)
    }

    function prepareStudentReqBody() {
        studentFormData["cource"] = parseInt(studentFormData["cource"])
        studentFormData["education_programme_id"] = parseInt(studentFormData["education_programme_id"])
        postReqOptions["body"] = JSON.stringify(studentFormData)
    }

    function RenderStudents() {
        return students.map((student) =>
            <option value={student.id}>
                {`${student.surname} ${student.name[0]}. ${student.middlename[0]}., ${student.cource} курс ${student.education_programme}`}
            </option>)
    }

    function RenderEdProgrammes() {
        return edprogs.map((edprog) =>
            <option value={edprog.id}>
                {`${edprog.name}`}
            </option>)
    }
    function OpenRequestResultModal() {
        setShowAddProjectesult(true);
    }
    function CloseRequestResultModal() {
        setShowAddProjectesult(false);
        setAddProjectResult(null);
    }

    function RenderRequestResultModal() {
        let header = "Научное руководство оформлено!";
        let body = "Вы можете просмотреть проект в списке проектов.";
        if (addProjectResult !== 200) {
            header = "Произошла ошибка при оформлении научного руководства!";
            body = `Код ошибки: ${addProjectResult}. Обратитесь в службу поддержки, если прблема не устранится.`;
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

    function moveNextStep(event, next, after) {
        let newArray = stepPath;
        if (next) {
            newArray = newArray.map((p, i) => {
                if (i === currentStepIndex + 1) {
                    return next;
                } else {
                    return p;
                }
            });
        }
        if (after) {
            newArray = newArray.map((p, i) => {
                if (i === currentStepIndex + 2) {
                    return after;
                } else {
                    return p;
                }
            });
        }
        setPath(newArray);
        setCurrentStep(currentStepIndex + 1);
    }
    function movePrevStep(event) {
        let newStep = methodStudentStep;
        if (currentStepIndex - 1 >= 0) {
            newStep = currentStepIndex - 1;
        }
        setCurrentStep(newStep);
    }


    function renderNavBar() {
        if (currentStepIndex === methodStudentStep) {
            setNavbar(renderStudentNavBar());
            return
        }
        setNavbar(renderPrevNextNavBar());
    }

    function renderStudentNavBar() {
        return <>
            <Nav as={Row} className="justify-content-between">
                <Nav.Item as={Col} sm={5} className='my-2'>
                    <Nav.Link className='border btn style-button link-body-emphasis'
                        onClick={(event) => { setAddingStudent(true); moveNextStep(event, addStudentStep, addProjectStep); }}
                        eventKey={steps[addStudentStep]}>Добавить нового студента</Nav.Link>
                </Nav.Item>
                <Nav.Item as={Col} sm={5} className='my-2'>
                    <Nav.Link className='border btn style-button link-body-emphasis'
                        onClick={(event) => { setAddingStudent(false); moveNextStep(event, selectStudentStep, addProjectStep) }}
                        eventKey={steps[selectStudentStep]}>Выбрать существующего студента</Nav.Link>
                </Nav.Item>
            </Nav>
        </>
    }

    function renderPrevNextNavBar() {
        let prevNav = currentStepIndex === 0 ? <></> :
            <Nav.Item as={Col} sm={5} className='my-2'>
                <Nav.Link className='border btn style-button-outline link-body-emphasis' onClick={movePrevStep} eventKey={steps[stepPath[currentStepIndex - 1]]}>Назад</Nav.Link>
            </Nav.Item>

        let nextNav = stepPath[currentStepIndex] + 1 === steps.length ? <></> :
            <Nav.Item as={Col} sm={5} className='my-2'>
                <Nav.Link disabled={addingStudent ? !isStudentFilled : !isStudentSelected} className='border btn style-button-outline link-body-emphasis'
                    onClick={(event) => currentStepIndex + 1 === stepPath.length ? moveNextStep(event, stepPath[currentStepIndex] + 1) : moveNextStep(event)}
                    eventKey={steps[stepPath[currentStepIndex + 1]]}>Далее</Nav.Link>
            </Nav.Item>

        return <>
            <Nav as={Row} className="justify-content-between" >
                {prevNav}
                {nextNav}
            </Nav>
        </>
    }

    return (
        <>
            <Row className='justify-content-center'>
                <Col xs={11} md={10} lg={8}>
                    <h1 className='mb-4'>Новое научное руководство</h1>
                    <hr />
                    {integr ?
                        integr.cloud_drive ?
                            <div ><Tab.Container id="left-tabs-example" defaultActiveKey="select-student-method">
                                <Row className='justify-content-center'>
                                    <Col xs={12} sm={8} className='mb-3'>
                                        <Tab.Content>
                                            <Tab.Pane eventKey="select-student-method">
                                                <div className='fs-3'>
                                                    Выберите студента, которого берете под научное руководство:
                                                </div>
                                            </Tab.Pane>
                                            <Tab.Pane eventKey="add-student">
                                                <div className='fs-3'>
                                                    Введите информацию о студенте:
                                                </div>
                                                <Form ref={addStudentRef} id="student-form">
                                                    <Form.Group className="mb-3">
                                                        <Form.Label>Имя *</Form.Label>
                                                        <Form.Control name="name" onChange={handleStudentChange} required placeholder="Введите название" />
                                                    </Form.Group>

                                                    <Form.Group className="mb-3" >
                                                        <Form.Label>Фамилия *</Form.Label>
                                                        <Form.Control name="surname" onChange={handleStudentChange} required placeholder="Введите описание" />
                                                    </Form.Group>

                                                    <Form.Group className="mb-3" >
                                                        <Form.Label>Отчество</Form.Label>
                                                        <Form.Control name="middlename" onChange={handleStudentChange} placeholder="Введите описание" />
                                                    </Form.Group>
                                                    <div className="mb-3" >
                                                        <label className='mb-2'>Курс *</label>
                                                        <select className="form-select" name="cource" onChange={handleStudentChange} required defaultValue="" >
                                                            <option value="" hidden>Выберете курс...</option>
                                                            <option value="1" >1</option>
                                                            <option value="2" >2</option>
                                                            <option value="3" >3</option>
                                                            <option value="4" >4</option>
                                                            <option value="5" >5</option>
                                                        </select>
                                                    </div>
                                                    <div className="mb-3">
                                                        <label className='mb-2'>Образовательная программа *</label>
                                                        <select className="form-select" name="education_programme_id" onChange={handleStudentChange} required defaultValue="">
                                                            <option value="" hidden>Выберете программу...</option>
                                                            {edprogs ? RenderEdProgrammes() : <option>загрузка...</option>}
                                                        </select>
                                                    </div>
                                                </Form>
                                            </Tab.Pane>
                                            <Tab.Pane eventKey="select-student">
                                                {students ?
                                                    <div className='fs-3'>
                                                        <Form.Group className="mb-3" controlId="student">
                                                            <Form.Label>Выберете студента:</Form.Label>
                                                            <Form.Select ref={selectRef} name="student_id" onChange={handleSelectChange} required aria-label="select student" defaultValue={-1} >
                                                                <option value={-1} hidden>Выберете студента...</option>
                                                                {RenderStudents()}
                                                            </Form.Select>
                                                        </Form.Group>
                                                    </div>
                                                    : SpinnerCenter()}

                                            </Tab.Pane>
                                            <Tab.Pane eventKey="add-project">
                                                <div className='fs-3'>
                                                    Введите информацию о проекте:
                                                    <Form id="project-form" onSubmit={handleSubmit}>
                                                        <Form.Group className="mb-3" controlId="meetName">
                                                            <Form.Label>Тема работы *</Form.Label>
                                                            <Form.Control name="theme" onChange={handleChange} required placeholder="Введите тему" />
                                                        </Form.Group>

                                                        <Form.Group className="mb-3" controlId="meetName">
                                                            <Form.Label>Год выполнения *</Form.Label>
                                                            <Form.Control type='number' name="year" onChange={handleChange} required placeholder="Введите год" />
                                                        </Form.Group>

                                                        <Form.Group className="mb-3" controlId="meetName">
                                                            <Form.Label>Логин обладателя репозитория *</Form.Label>
                                                            <Form.Control name="repository_owner_login" onChange={handleChange} required placeholder="Введите логин" />
                                                        </Form.Group>

                                                        <Form.Group className="mb-3" controlId="meetName">
                                                            <Form.Label>Имя репозитория *</Form.Label>
                                                            <Form.Control name="repository_name" onChange={handleChange} required placeholder="Введите название" />
                                                        </Form.Group>
                                                        <Row className='justify-content-center mx-1'>
                                                            <Button type="submit" className="style-button mb-3">
                                                                Взять под научное руководство
                                                            </Button>
                                                        </Row>
                                                    </Form>
                                                </div>
                                            </Tab.Pane>
                                        </Tab.Content>
                                    </Col>
                                    <Col xs={12} sm={8}>
                                        {navbar}
                                    </Col>
                                </Row>
                            </Tab.Container>
                            </div>
                            :
                            <>
                                <h3>Вы еще не подключили облачное хранилище, это можно сделать <a href='/profile/integrations'>здесь</a></h3>
                            </>
                        : SpinnerCenter()}



                </Col>
            </Row>

            <Modal show={showAddProjectResult} onHide={CloseRequestResultModal}>
                {addProjectResult ? RenderRequestResultModal() :
                    <Modal.Header className="justify-content-md-center">
                        {SpinnerCenter()}
                    </Modal.Header>}
            </Modal>
        </>
    );
}

export default AddProject;