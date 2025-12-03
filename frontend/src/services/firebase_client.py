import uuid

class FirebaseClient:
    def __init__(self):
        self.user_id = None
    
    def sign_in_anonymously(self) -> str:
        if not self.user_id:
            self.user_id = f"anonymous-{uuid.uuid4().hex[:8]}"
        return self.user_id
    
    def get_current_user_id(self) -> str:
        if not self.user_id:
            return self.sign_in_anonymously()
        return self.user_id

