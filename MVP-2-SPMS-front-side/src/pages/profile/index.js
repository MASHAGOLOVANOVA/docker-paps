import React, { useState, useEffect } from 'react';
import Cookies from 'universal-cookie';
import { Col, Row, } from 'react-bootstrap';
import SpinnerCenter from 'pages/shared/Spinner';
import ProfileSidebar from './ProfileSidebar';

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

function Profile() {
    const [user, setUser] = useState(null);


    useEffect(() => {
        fetch(`${process.env.REACT_APP_SERVER_ADDR}/api/v1/account`, reqOptions)
            .then(response => response.json())
            .then(json => setUser(json))
            .catch(error => console.error(error));
    }, []);
    return (
        <>
            <Row className='m-2'>
                <Col xs={12} sm={12} md={4} lg={2}>
                    <ProfileSidebar/>
                </Col>
                <Col xs={12} sm={12} md={8} lg={10} className='px-5'>
                    {user ? <>
                        <h3 className='mb-2'>{user.name}</h3>
                        <div className='fs-5 mb-2 fst-italic'>{user.science_degree}</div>
                        <hr />
                        <div>
                            <Row className='' >
                                <Col >
                                    <div className='fs-3 mb-2 fw-medium'>{user.university}</div>
                                </Col>
                            </Row>
                            <Row className='mb-3' xs="auto">
                                <Col >
                                    <div className='fs-4 mb-2 fw-light'>Свободные места:</div>                                    
                                </Col>
                                <Col >
                                {SpinnerCenter()}                                  
                                </Col>
                            </Row>
                        </div>
                    </> : SpinnerCenter()}
                </Col>
            </Row>
        </>
    );
};

export default Profile;