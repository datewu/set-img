import { useState, useEffect } from 'react'
import reactLogo from './assets/react.svg'
import './App.css'

import CURL from './curl';

function App() {
  const [data, setData] = useState(null);
  useEffect(() => {
    const getData = async () => {
      const resp = await fetch("/api/v1");
      if (!resp.ok) {
        throw resp.statusText
      }
      const json = await resp.json();
      console.log("json", json);
      setData(json)
      return json;
    }
    getData()
  }, [])

  return (
    <div className="App">
      <div>
        <a href="https://github.com/datewu/set-img" target="_blank">
          <img src="/fav.svg" className="logo" alt="logo" />
        </a>
        <a href="https://reactjs.org" target="_blank">
          <img src={reactLogo} className="logo react" alt="React logo" />
        </a>
      </div>
      <h1>kubernetes部署工具</h1>
      <div className="card">
        <CURL data={data} />
      </div>
      <p className="read-the-docs">
        Click on the set-img and React logos to learn more
      </p>
    </div>
  )
}

export default App
