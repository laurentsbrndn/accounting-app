import 'package:flutter/material.dart';
import '../../widgets/login/login_form.dart';

class LoginPage extends StatelessWidget {
  const LoginPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SafeArea(
        child: Center(
          child: SingleChildScrollView(
            padding: const EdgeInsets.symmetric(horizontal: 24.0),
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                // Logo placeholder (ganti dengan Image.asset/logo nanti)
                const Icon(Icons.ac_unit, size: 80, color: Color(0xFF3F51B5)),

                const SizedBox(height: 24),
                const Text(
                  'Log In Now',
                  style: TextStyle(
                    fontSize: 24,
                    fontWeight: FontWeight.bold,
                    color: Color(0xFF3F51B5),
                  ),
                ),
                const SizedBox(height: 8),
                const Text(
                  'Please login to continue using our app',
                  style: TextStyle(color: Colors.black54),
                ),

                const SizedBox(height: 32),
                const LoginForm(),

                const SizedBox(height: 16),
                Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    const Text("Don’t have an account? "),
                    GestureDetector(
                      onTap: () {
                        // Arahkan ke halaman sign up
                      },
                      child: const Text(
                        "Sign Up",
                        style: TextStyle(
                          color: Colors.blue,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    )
                  ],
                ),

                const SizedBox(height: 24),
                const Text("Or connect with"),
                const SizedBox(height: 12),
                Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    socialButton(icon: Icons.facebook),
                    const SizedBox(width: 16),
                    // socialButton(icon: Icons.twitter),
                    // const SizedBox(width: 16),
                    // socialButton(icon: Icons.linkedin),
                  ],
                )
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget socialButton({required IconData icon}) {
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        shape: BoxShape.circle,
        border: Border.all(color: Colors.grey.shade300),
      ),
      child: Icon(icon, color: Color(0xFF1E88E5)),
    );
  }
}