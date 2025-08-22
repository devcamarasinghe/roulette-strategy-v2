from fastapi import FastAPI, HTTPException, Depends
from pydantic import BaseModel
from sqlalchemy import create_engine, Column, Integer, String, DateTime
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker, Session
from datetime import datetime

# Database configuration
DATABASE_URL = "postgresql://roulette_user:roulette_password@postgres-db:5432/roulette_db"

# SQLAlchemy setup
engine = create_engine(DATABASE_URL)
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)
Base = declarative_base()

# Database model
class UserDB(Base):
    __tablename__ = "users"
    
    id = Column(Integer, primary_key=True, index=True)
    username = Column(String, unique=True, index=True)
    email = Column(String, unique=True, index=True)
    password = Column(String)  # In production, hash this!
    created_at = Column(DateTime, default=datetime.utcnow)

# Create tables
Base.metadata.create_all(bind=engine)

# Create the "Customer Service Building"
app = FastAPI(title="User Service", version="2.0.0 - PostgreSQL Edition")

# Database dependency
def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()

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
async def register_user(user: UserRegister, db: Session = Depends(get_db)):
    """Register a new customer - now with PostgreSQL!"""
    
    # Check if user already exists
    db_user = db.query(UserDB).filter(UserDB.username == user.username).first()
    if db_user:
        raise HTTPException(
            status_code=400, 
            detail="Username already registered"
        )
    
    # Check if email already exists  
    db_email = db.query(UserDB).filter(UserDB.email == user.email).first()
    if db_email:
        raise HTTPException(
            status_code=400,
            detail="Email already registered"
        )
    
    # Create new user in database
    db_user = UserDB(
        username=user.username,
        email=user.email,
        password=user.password  # In real app, hash this!
    )
    db.add(db_user)
    db.commit()
    db.refresh(db_user)
    
    return UserResponse(
        username=db_user.username,
        email=db_user.email,
        message=f"Welcome {db_user.username}! Registration successful in PostgreSQL!"
    )

@app.post("/login")
async def login_user(credentials: UserLogin, db: Session = Depends(get_db)):
    """Customer login - now with PostgreSQL verification!"""
    
    # Check if user exists
    db_user = db.query(UserDB).filter(UserDB.username == credentials.username).first()
    if not db_user:
        raise HTTPException(
            status_code=401,
            detail="User not found"
        )
    
    # Check password (simple comparison for learning)
    if db_user.password != credentials.password:
        raise HTTPException(
            status_code=401,
            detail="Invalid password"
        )
    
    return {
        "message": f"Welcome back, {db_user.username}! (from PostgreSQL)",
        "status": "login_successful",
        "user": {
            "username": db_user.username,
            "email": db_user.email,
            "created_at": db_user.created_at.isoformat()
        }
    }

@app.get("/users")
async def get_users(db: Session = Depends(get_db)):
    """List all registered users - from PostgreSQL database!"""
    
    # Get all users from database
    users = db.query(UserDB).all()
    
    if not users:
        return {"message": "No users registered yet", "source": "PostgreSQL"}
    
    # Return users without passwords
    users_list = []
    for user in users:
        users_list.append({
            "username": user.username,
            "email": user.email,
            "created_at": user.created_at.isoformat()
        })
    
    return {
        "total_users": len(users_list),
        "users": users_list,
        "source": "PostgreSQL Database ✅"
    }

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8001)
