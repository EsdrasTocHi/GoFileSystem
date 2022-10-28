import React, {useRef} from "react";
import Col from "react-bootstrap/esm/Col"
import Form from 'react-bootstrap/esm/Form';
import Container from 'react-bootstrap/esm/Container'
import Row from 'react-bootstrap/esm/Row'
import Button from 'react-bootstrap/esm/Button'
import FloatingLabel from "react-bootstrap/esm/FloatingLabel";
import 'bootstrap/dist/css/bootstrap.min.css';
import './Consoles.css'

class Console extends React.Component{
    constructor(props){
        super(props)

        this.state ={
            inputRef : useRef<HTMLInputElement>(null)
        }
    }

    handleFile(){
        let ref = (document.getElementById("file") as HTMLInputElement)
        this.setState(
            {
                inputRef : 
            }
        )
    }

    render(){
        return(
            <Col xs={12} md={6}>
                <Container>
                    <Row>
                        <Col xs={12} md={12}>
                            <Form.Group id="fileChooser" controlId="formFile" className="mb-3">
                                <Form.Control id="file" type="file" accept=".mia"/>
                            </Form.Group>
                            <Button variant="primary" id="btnExecute">Execute</Button>
                        </Col>
                    </Row>
                    <Row>
                        <Col xs={12} md={12}>
                            <FloatingLabel controlId="floatingTextarea2" label="Commands">
                                <Form.Control
                                as="textarea"
                                id="txtCommand"
                                readOnly
                                ></Form.Control>
                            </FloatingLabel>
                        </Col>
                    </Row>
                    <Row>
                        <Col xs={12} md={12}>
                            <FloatingLabel controlId="floatingTextarea2" label="Console">
                                <Form.Control
                                as="textarea"
                                id="txtRes"
                                readOnly
                                ></Form.Control>
                            </FloatingLabel>
                        </Col>
                    </Row>
                </Container>
            </Col>
        );
    }
}

export default Console