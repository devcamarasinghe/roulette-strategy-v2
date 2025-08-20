from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import Dict

# Create the "Customer Service Building"
app = FastAPI(title="User Service", version="1.0.0")

# Simple in-memory user database (for learning purposes)
users_db: Dict[str, Dict] = {}

# Pydantic models for request/response data
class UserRegister(BaseModel):
    username: str
    email: str
    password: str

class UserLogin(BaseModel):
    username: str
    password: str

class UserResponse(BaseModel):
    username: str
    email: str
    message: str

@app.get("/")
async def welcome():
    """Welcome message - like the reception desk"""
    return {
        "message": "Welcome to User Service!", 
        "status": "healthy",
        "service": "user-service",
        "port": 8001,
        "endpoints": {
            "register": "POST /register",
            "login": "POST /login",
            "users": "GET /users"
        }
    }

@app.get("/health")
async def health_check():
    """Health check - is our service working?"""
    return {"status": "healthy", "service": "Customer Service Building"}

@app.post("/register", response_model=UserResponse)
async def register_user(user: UserRegister):
    """Register a new customer - like opening an account"""
    
    # Check if user already exists
    if user.username in users_db:
        raise HTTPException(
            status_code=400, 
            detail="Username already registered"
        )
    
    # Store user in our simple database
    users_db[user.username] = {
        "username": user.username,
        "email": user.email,
        "password": user.password,  # In real app, hash this!
        "created_at": "2025-08-20"
    }
    
    return UserResponse(
        username=user.username,
        email=user.email,
        message=f"Welcome {user.username}! Registration successful."
    )

@app.post("/login")
async def login_user(credentials: UserLogin):
    """Customer login - like showing ID card"""
    
    # Check if user exists
    if credentials.username not in users_db:
        raise HTTPException(
            status_code=401,
            detail="User not found"
        )
    
    # Check password (simple comparison for learning)
    stored_user = users_db[credentials.username]
    if stored_user["password"] != credentials.password:
        raise HTTPException(
            status_code=401,
            detail="Invalid password"
        )
    
    return {
        "message": f"Welcome back, {credentials.username}!",
        "status": "login_successful",
        "user": {
            "username": stored_user["username"],
            "email": stored_user["email"]
        }
    }

@app.get("/users")
async def get_users():
    """List all registered users - like customer directory"""
    if not users_db:
        return {"message": "No users registered yet"}
    
    # Return users without passwords
    users_list = []
    for username, user_data in users_db.items():
        users_list.append({
            "username": user_data["username"],
            "email": user_data["email"],
            "created_at": user_data["created_at"]
        })
    
    return {
        "total_users": len(users_list),
        "users": users_list
    }

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8001)
