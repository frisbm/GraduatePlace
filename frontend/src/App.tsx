import {createFetch} from "./store";

const App = () => {
    const result = createFetch<object>("/user", "GET", true);

    return (
        <>
            <h1>{JSON.stringify(result.data)}</h1>
            <h1>{JSON.stringify(result.loading)}</h1>
            <h1>{JSON.stringify(result.error)}</h1>
        </>
    );
}

export default App;
