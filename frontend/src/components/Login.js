import React from 'react'
import Container from 'react-bootstrap/Container'
import FormGroup from 'react-bootstrap/esm/FormGroup'
import Row from 'react-bootstrap/esm/Row'
import Button from 'react-bootstrap/esm/Button'

class Login extends React.Component{
    constructor(props){
        super(props)

        this.state = {
            user : "",
            passw : "",
            id : "",
            active : false
        }

        this.path = "http://3.144.197.243:3030"

        this.handleInputChange.bind(this.handleInputChange)
        this.login.bind(this.login)
        this.logout.bind(this.logout)
    }

    handleInputChange(event){
        const target = event.target;
        const name = target.name;
        const value = target.value;
        
        this.setState({
            [name]: value
        });
    }

    login(){
        const comm = {
            command : "login -usuario="+this.state.user+" -password="+this.state.passw+" -id="+this.state.id
        }

        let requestPost = {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(comm),
          };
          const url = this.path+'/command';
        fetch(url, requestPost)
            .then((response) => response.json())
            .then((data) => {
                console.log(data)
                if(!data.error){
                    this.setState({active:true})
                }

                alert(data.response)
            })
    }

    logout(){
        const comm = {
            command : "logout"
        }

        let requestPost = {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(comm),
          };
          const url = this.path+'/command';
        fetch(url, requestPost)
            .then((response) => response.json())
            .then((data) => {
                if(!data.error){
                    this.setState({active:false})
                }

                alert(data.response)
            })
    }

    render(){
        if (!this.state.active){
        return(
            <Container>
                <Row>
                    <form>
                        <FormGroup>
                            <label htmlFor="txtUser">User</label>
                            <input type="text" name='user' className="form-control" id="txtUser" onChange={(e)=>this.handleInputChange(e)}/>
                        </FormGroup>
                        <FormGroup>
                            <label htmlFor="txtPassw">Password</label>
                            <input type="password" name='passw' className="form-control" id="txtPassw" onChange={(e)=>this.handleInputChange(e)}/>
                        </FormGroup>
                        <FormGroup>
                            <label htmlFor="txtId">Password</label>
                            <input type="text" name='id' className="form-control" id="txtId" onChange={(e)=>this.handleInputChange(e)}/>
                        </FormGroup>
                        <FormGroup>
                            <Button variant="primary" id="btnLogin" onClick={()=>this.login()}>Login</Button>
                        </FormGroup>
                    </form>
                </Row>
            </Container>
        )
        }else{
            return(
                <Container>
                    <Row>
                        <form>
                            <FormGroup>
                                <label htmlFor="txtUser">User</label>
                                <input type="text" className="form-control" id="txtUser" readOnly/>
                            </FormGroup>
                            <FormGroup>
                                <label htmlFor="txtPassw">Password</label>
                                <input type="password" className="form-control" id="txtPassw" readOnly/>
                            </FormGroup>
                            <FormGroup>
                                <label htmlFor="txtId">Password</label>
                                <input type="text" className="form-control" id="txtId" readOnly/>
                            </FormGroup>
                            <FormGroup>
                                <Button variant="primary" id="btnLogout" onClick={()=>this.logout()}>Logout</Button>
                            </FormGroup>
                        </form>
                    </Row>
                </Container>
            )
        }
    }
}

export default Login