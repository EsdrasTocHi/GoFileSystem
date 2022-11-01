import './App.css';
import Container from 'react-bootstrap/esm/Container';
import Col from 'react-bootstrap/esm/Col';
import Console from './components/Consoles.js'
import Login from './components/Login';
import Row from 'react-bootstrap/esm/Row';

function App() {
  return (
    <>
    <header>
      <h1>
        FileSystem
      </h1>
    </header>
    <Container id='generalContainer'>
      <Row>
      <Console></Console>
      <Col xs={12} md={6}>
        <Login></Login>
      </Col>
      </Row>
    </Container>
    </>
  );
}

export default App;
