package com.inszva.monkey;

import java.io.IOException;
import java.net.UnknownHostException;

class main{
	public static void main(String[] args) throws UnknownHostException, IOException {
		Monkey monkey = new Monkey("127.0.0.1",1517,"monkey");
		for(int i = 0;i < 100;i++) {
			monkey.set(String.valueOf(i), String.valueOf(i));
		}
		monkey.release();
	}
}
