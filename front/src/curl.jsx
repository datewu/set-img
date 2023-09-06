import { useState } from 'react'

function CURL({ data }) {
  const [detail, showDetail] = useState(false);
  const [copied, setCopied] = useState(false);
  if (!data) {
    return null
  }
  const cp = () => {
    navigator.clipboard.writeText(data.url);
    setCopied(true)
  }

  let pre;
  if (detail) {
    pre = <pre style={{ width: 550, whiteSpace: 'pre-wrap', textAlign: 'left' }}>
      curl -X{data.method} {data.url} \ <br />
      -H 'Authorization: $TOKEN' \ <br />
      --data-binary  \ <br />
      <code>
        {JSON.stringify(data.data)}
      </code>
    </pre>
  }

  return (
    <div>
      <div className='url-btns'>
        <p>
          {data.url}
        </p>
        <div>
          <button onClick={cp}>{copied ? 'copied' : 'click to copy'} </button>
          <button onClick={() => showDetail(!detail)}>{detail ? 'less' : 'more'}</button>
        </div>
      </div>
      <div>
        {pre}
      </div>
    </div >
  )
}

export default CURL;
