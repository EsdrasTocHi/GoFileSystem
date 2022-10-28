import './App.css';
import Container from 'react-bootstrap/esm/Container';
import Console from './components/Consoles.js'

function App() {
  return (
    <>
    <header>
      <h1>
        FileSystem
      </h1>
    </header>
    <Container id='generalContainer'>
      <Console></Console>
    </Container>
    </>
  );
}

export default App;
