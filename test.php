<?php
class MonkeyDB
{
	private $socket;
	
	private function read()
	{
		$data = "";
		$total = 0;
		$t = fread($this->socket,4);
		for($i = 0;$i < 4;$i++)
		{
			$total *= 256;
			$total += ord($t[$i]);
		}
		while($total > 1024)
		{
			$buf = fread($this->socket,1024);
			$data .= $buf;
			$total-=1024;
		}
		$buf = fread($this->socket,$total);
		$data .= $buf;
		return $data;
	}
	
	private function write($string)
	{
		$total = strlen($string);
		fwrite($this->socket,strrev(pack("L",$total)));
		for($i = 0;$i < $total - 1024;$i += 1024)
		{
			fwrite($this->socket,substr($string,$i * 1024,1024));
		}
		fwrite($this->socket,substr($string,$i * 1024));
	}
	
	public function __construct($addr,$passwd)
	{
		set_time_limit(0);
		ob_implicit_flush();
		$this->socket = stream_socket_client("tcp://{$addr}:1517");
		if(!$this->socket)
		{
			throw new Exception("monkey 连接失败！");
		}
		$this->write("auth " . $passwd);
		$this->read();
	}
	
	public function __destruct()
	{
		fclose($this->socket);
	}
	
	//data: string	serialize if necessary
	public function set($key,$data)
	{
		if(!$this->socket)
		{
			throw new Exception("monkey 连接失败！");
		}
		$this->write("set {$key} ".($data));
		$this->read();
	}
	
	//return: string
	public function get($key)
	{
		if(!$this->socket)
		{
			throw new Exception("monkey 连接失败！");
		}
		$this->write("get {$key}");
		$data = $this->read();
		return ($data);
	}
	
	public function remove($key)
	{
		if(!$this->socket)
		{
			throw new Exception("monkey 连接失败！");
		}
		$this->write("remove {$key}");
		$this->read();
	}
	
	public function createDB($dbName)
	{
		if(!$this->socket)
		{
			throw new Exception("monkey 连接失败！");
		}
		$this->write("createdb {$dbName}");
		$this->read();
	}
	
	public function switchDB($dbName)
	{
		if(!$this->socket)
		{
			throw new Exception("monkey 连接失败！");
		}
		$this->write("switchdb {$dbName}");
		$data = $this->read();
	}
	
	public function dropDB($dbName)
	{
		if(!$this->socket)
		{
			throw new Exception("monkey 连接失败！");
		}
		$this->write("dropdb {$dbName}");
	}
	
	//return string
	public function listDB()
	{
		if(!$this->socket)
		{
			throw new Exception("monkey 连接失败！");
		}
		$this->write("listdb ");
		$data = $this->read();
		return $data;
	}
}

//  $monkey = new MonkeyDB("127.0.0.1","monkey");
//   for($i = 0;$i < 100;$i++)
//   	$monkey->get("{$i}");
