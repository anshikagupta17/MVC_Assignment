import RegisterForm from "./RegisterForm";

export default function RegisterPage() {
    return (
        <div
            className="min-h-screen flex items-center justify-center"
            style={{ backgroundImage: `url('/assets/battle.png')` }}
        >
            <div className="bg-[white]/60 text-black p-8 rounded-xl shadow-lg w-full max-w-sm">
                <h1 className="text-3xl font-bold mb-6 text-center">
                    Register
                </h1>

                <RegisterForm />
            </div>
        </div>
    );
}