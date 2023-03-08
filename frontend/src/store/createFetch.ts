import {createResource, createSignal} from "solid-js";
import {createStore, SetStoreFunction} from "solid-js/store";


type Tokens = {
    accessToken?: string
    refreshToken?: string
}

const defaultTokens = {
    accessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzgyMzAzMjAsImlhdCI6MTY3ODIzMDMxOSwia2luZCI6ImFjY2VzcyIsInV1aWQiOiJhMjc2M2EyYS03ZjY1LTRiZjQtYmMwYy1lZGVjY2Q1MjI2YjUifQ.glJvDcOUiGCMVV-w8Mhtf3UUcq20tWlWxg_3l1veynU",
    refreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODA4MjIzMTksImlhdCI6MTY3ODIzMDMxOSwia2luZCI6InJlZnJlc2giLCJ1dWlkIjoiYTI3NjNhMmEtN2Y2NS00YmY0LWJjMGMtZWRlY2NkNTIyNmI1In0.Fkp7z4bGna4ZwO7BkAu5vetRPk-s5QUo63x1C9A4pXk",
}

type MethodType = "GET" | "POST"

type CreateFetchResult<T> = {
    data: T | null,
    loading: boolean,
    error: string | null,
}

const fetchNoAuth = async <T>(result: CreateFetchResult<T>, setResult: SetStoreFunction<CreateFetchResult<T>>, endpoint: string, method: MethodType, body?: any): Promise<T> => {
    setResult({
        ...result,
        loading: true,
    })
    const request = fetch(
        'http://localhost:8080/api' + endpoint,
        {
            method: method,
            headers: {
                "Content-Type": "application/json",
            },
            body: body ? JSON.stringify(body) : null,
        },
    )
    const response = await request;
    if (response.status >= 400) {
        setResult({
            ...result,
            loading: false,
            error: await response.json(),
        });
    } else {
        setResult({
            ...result,
            loading: false,
            data: await response.json(),
        });
    }

    return response.json();
};

const fetchAuth = async <T>(result: CreateFetchResult<T>, setResult: SetStoreFunction<CreateFetchResult<T>>, endpoint: string, method: MethodType, body?: any): Promise<T> => {
    localStorage.setItem("tokens", JSON.stringify(defaultTokens))
    const tokens = JSON.parse(`${localStorage.getItem("tokens")}`) as Tokens

    setResult({
        ...result,
        loading: true,
    })

    const request = fetch(
        'http://localhost:8080/api' + endpoint,
        {
            method: method,
            headers: {
                "Content-Type": "application/json",
                "Authorization": `${tokens.accessToken}`,
            },
            body: body ? JSON.stringify(body) : null,
        },
    )
    const response = await request;
    if (response.status >= 400) {
        const error = await response.json();
        if (error == "access token expired") {
            let response = await fetch(
                'http://localhost:8080/api/user/refresh',
                {
                    method: 'POST',
                    body: JSON.stringify({
                        refreshToken: tokens.refreshToken,
                    }),
                },
            )
            if (response.status >= 400) {
                setResult({
                    data: null,
                    error: await response.json() as string,
                    loading: false,
                })
                return response.json()
            }
            const newTokens = await response.json() as Tokens;
            localStorage.setItem("tokens", JSON.stringify({
                accessToken: newTokens.accessToken,
                refreshToken: newTokens.refreshToken,
            }));
            const request = fetch(
                'http://localhost:8080/api' + endpoint,
                {
                    method: method,
                    headers: {
                        "Content-Type": "application/json",
                        "Authorization": `${newTokens.accessToken}`,
                    },
                    body: body ? JSON.stringify(body) : null,
                },
            )
            response = await request;

            setResult({
                error: null,
                loading: false,
                data: await response.json(),
            });
            return response.json()

        }

    }
    return response.json()
};

const createFetch = <T>(endpoint: string, method: MethodType, auth?: boolean, body?: any): CreateFetchResult<T> => {
    const [_endpoint, setEndpoint] = createSignal<string>(endpoint)
    const [result, setResult] = createStore<CreateFetchResult<T>>({
        data: null,
        loading: false,
        error: null,
    });


    if (!auth) {
        createResource(_endpoint, async () => fetchNoAuth(result, setResult, endpoint, method, body))
        return result
    }

    createResource(_endpoint, async () => fetchAuth(result, setResult, endpoint, method, body));
    return result
}

export default createFetch;
