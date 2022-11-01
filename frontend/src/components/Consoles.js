import React from "react";
import Col from "react-bootstrap/esm/Col"
import Form from 'react-bootstrap/esm/Form';
import Container from 'react-bootstrap/esm/Container'
import Row from 'react-bootstrap/esm/Row'
import Button from 'react-bootstrap/esm/Button'
import FloatingLabel from "react-bootstrap/esm/FloatingLabel";
import 'bootstrap/dist/css/bootstrap.min.css';
import './Consoles.css'


function removeComment(line) {
    let res = ""

    for(let i = 0; i < line.length; i++){
        if (line[i] === '#') {
            break
        }

        res += line[i]
    }

    return res
}

class Console extends React.Component{
    constructor(props){
        super(props)

        this.state = {
            file : "",
            console : "",
            commands : ""
        }

        this.path = "http://127.0.0.1:3030"

        this.showFile.bind(this.showFile)
        this.execute.bind(this.execute)
    }

    async showFile(event){
        event.preventDefault()
        const reader = new FileReader()
        reader.onload = async (event) => { 
            const text = (event.target.result)
            alert(text)
            this.setState({
                file : text
            });
        };
        reader.readAsText(event.target.files[0])
    }

    async execute(){
        let lines = this.state.file.split("\n")
        let finalContent = ""
        let cons = ""

        for(let i = 0; i < lines.length; i++){
            let line = lines[i]
            finalContent += line+"\n"
            line = removeComment(line)
            line = line.trim()

            if (line !== ""){
                const comm = {
                    command : line
                }
    
                let requestPost = {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify(comm),
                  };
                  const url = this.path+'/command';
                 await fetch(url, requestPost)
                    .then((response) => response.json())
                    .then((data) => {
                        console.log(data);
                        cons += data.response +"\n";
                    })
            }
        }

        this.setState({
            commands:finalContent,
            console : cons
        })
    }

    render(){
        return(
            <Col xs={12} md={6}>
                <Container>
                    <Row>
                        <Col xs={12} md={12}>
                            <Form.Group id="fileChooser" className="mb-3">
                                <Form.Control id="file" type="file" accept=".script" onChange={(event) => this.showFile(event)}/>
                            </Form.Group>
                            <Button variant="primary" id="btnExecute" onClick={() => this.execute()}>Execute</Button>
                        </Col>
                    </Row>
                    <Row>
                        <Col xs={12} md={12}>
                            <FloatingLabel id="floatingTextarea2" label="Commands">
                                <Form.Control
                                as="textarea"
                                id="txtCommand"
                                value={this.state.commands}
                                readOnly
                                ></Form.Control>
                            </FloatingLabel>
                        </Col>
                    </Row>
                    <Row>
                        <Col xs={12} md={12}>
                            <FloatingLabel id="floatingTextarea2" label="Console">
                                <Form.Control
                                as="textarea"
                                id="txtRes"
                                value={this.state.console}
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