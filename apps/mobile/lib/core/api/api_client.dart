import 'dart:convert';
import 'dart:io';
import 'package:flutter/foundation.dart';
import 'package:http/http.dart' as http;
import 'package:shared_preferences/shared_preferences.dart';

class ApiClient {
  static final ApiClient _instance = ApiClient._internal();
  factory ApiClient() => _instance;
  ApiClient._internal();

  String get baseUrl {
    if (kIsWeb) return 'http://127.0.0.1:8080/api/v1';
    if (Platform.isAndroid) return 'http://10.0.2.2:8080/api/v1';
    return 'http://127.0.0.1:8080/api/v1';
  }

  String? _token;

  Future<void> init() async {
    final prefs = await SharedPreferences.getInstance();
    _token = prefs.getString('bn_token');
  }

  Future<void> setToken(String token) async {
    _token = token;
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString('bn_token', token);
  }

  Future<void> clearToken() async {
    _token = null;
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove('bn_token');
  }

  bool get isAuthenticated => _token != null && _token!.isNotEmpty;

  Map<String, String> get _headers {
    final headers = {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
    };
    if (_token != null) {
      headers['Authorization'] = 'Bearer $_token';
    }
    return headers;
  }

  Future<dynamic> get(String endpoint) async {
    final uri = Uri.parse('$baseUrl$endpoint');
    final resp = await http.get(uri, headers: _headers);
    return _handleResponse(resp);
  }

  Future<dynamic> post(String endpoint, Map<String, dynamic> body) async {
    final uri = Uri.parse('$baseUrl$endpoint');
    final resp = await http.post(uri, headers: _headers, body: jsonEncode(body));
    return _handleResponse(resp);
  }

  Future<dynamic> patch(String endpoint, Map<String, dynamic> body) async {
    final uri = Uri.parse('$baseUrl$endpoint');
    final resp = await http.patch(uri, headers: _headers, body: jsonEncode(body));
    return _handleResponse(resp);
  }

  Future<dynamic> delete(String endpoint) async {
    final uri = Uri.parse('$baseUrl$endpoint');
    final resp = await http.delete(uri, headers: _headers);
    return _handleResponse(resp);
  }

  dynamic _handleResponse(http.Response resp) {
    if (resp.body.isEmpty) return null;
    final json = jsonDecode(utf8.decode(resp.bodyBytes));
    if (resp.statusCode >= 200 && resp.statusCode < 300) {
      return json['data'];
    }
    final error = json['error'];
    throw ApiException(
      code: error?['code'] ?? 'UNKNOWN_ERROR',
      message: error?['message'] ?? 'An error occurred',
    );
  }
}

class ApiException implements Exception {
  final String code;
  final String message;
  ApiException({required this.code, required this.message});

  @override
  String toString() => message;
}
