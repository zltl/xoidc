
const PAGE_SIZE = 20;

export default async function OIDCClientPage({
  params: {
    lng,
    page
  }
}: {
  params: {
    lng: string
    page: string[]
  }
}) {
  let curPage = 1;
  if (page && page.length > 0) {
    curPage = parseInt(page[0]);
  }
  if (curPage < 1) {
    curPage = 1;
  }

  const data = await fetchClientList((curPage - 1) * PAGE_SIZE, PAGE_SIZE);
  console.log("data=", JSON.stringify(data));
  if (data.status != 'success') {
    return (
      <div>
        {data.msg}
      </div>
    )
  }

  const clientsList = data.clients.map((client: any) => {
    return (
      <div key={client.id}>
        <div>{client.id}</div>
        <div>{client.secret}</div>
      </div>
    )
  });

  return (
    <div>
      <div>
        {clientsList}
      </div>
    </div>
  )
}

async function fetchClientList(offset: number, limit: number) {
  const res = await fetch(process.env.API_URL + `/api/client?offset=${offset}&limit=${limit}`);
  const data = await res.json();
  return data;
}

